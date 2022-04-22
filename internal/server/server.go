// Package server provides the bootstrapping logic for a node in an arrebato cluster.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/hashicorp/serf/serf"
	"go.etcd.io/bbolt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/davidsbond/arrebato/internal/acl"
	"github.com/davidsbond/arrebato/internal/backup"
	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/consumer"
	"github.com/davidsbond/arrebato/internal/message"
	"github.com/davidsbond/arrebato/internal/node"
	"github.com/davidsbond/arrebato/internal/prune"
	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/topic"
)

type (
	// The Server type represents the entire arrebato node and stitches together the serf, raft & grpc configuration.
	Server struct {
		config  Config
		logger  hclog.Logger
		store   *bbolt.DB
		restore chan struct{}
		pruner  *prune.Pruner

		// Dependencies for Serf
		serf       *serf.Serf
		serfEvents <-chan serf.Event

		// Dependencies for Raft
		raft          *raft.Raft
		raftStore     *raftboltdb.BoltStore
		raftTransport *transport.Manager

		// Dependencies for ACLs
		aclStore   *acl.BoltStore
		aclHandler *acl.Handler
		aclGRPC    *acl.GRPC

		// Dependencies for Topics
		topicStore   *topic.BoltStore
		topicHandler *topic.Handler
		topicGRPC    *topic.GRPC

		// Dependencies for Messages
		messageStore   *message.BoltStore
		messageHandler *message.Handler
		messageGRPC    *message.GRPC

		// Dependencies for Consumers
		consumerStore   *consumer.BoltStore
		consumerHandler *consumer.Handler
		consumerGRPC    *consumer.GRPC

		// Dependencies for Signing Keys
		signingStore   *signing.BoltStore
		signingHandler *signing.Handler
		signingGRPC    *signing.GRPC

		// Dependencies for Nodes
		nodeGRPC *node.GRPC
	}

	// The Config type contains configuration values for the Server.
	Config struct {
		// The server version
		Version string

		// LogLevel denotes the verbosity of logs.
		LogLevel int

		// BindAddress denotes the address the server will bind to for serf, raft & grpc.
		BindAddress string

		// AdvertiseAddress denotes the address the server will advertise to other nodes for serf & raft.
		AdvertiseAddress string

		// DataPath denotes the on-disk location that raft & message state are stored.
		DataPath string

		// Peers contain existing node addresses that should be connected to on start.
		Peers []string

		// PruneInterval determines how frequently messages are pruned from topics based on the topic's
		// retention period.
		PruneInterval time.Duration

		Raft    RaftConfig
		Serf    SerfConfig
		GRPC    GRPCConfig
		TLS     TLSConfig
		Metrics MetricConfig
	}

	// The TLSConfig type contains configuration values for TLS between the server nodes and clients.
	TLSConfig struct {
		// Location of the TLS certificate file to use for transport credentials.
		CertFile string

		// Location of the TLS Key file to use for transport credentials.
		KeyFile string

		// The certificate of the CA that signs client certificates.
		CAFile string
	}
)

// New returns a new instance of the Server type based on the provided Config.
func New(config Config) (*Server, error) {
	var err error
	server := &Server{
		restore: make(chan struct{}, 1),
		logger: hclog.New(&hclog.LoggerOptions{
			Level:  hclog.Level(config.LogLevel),
			Output: os.Stdout,
		}),
	}

	if err = os.MkdirAll(config.DataPath, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create data path: %w", err)
	}

	if config.AdvertiseAddress == "" {
		config.AdvertiseAddress, err = getDefaultAdvertiseAddress()
		if err != nil {
			return nil, err
		}

		server.logger.Info("using default advertise address", "ip", config.AdvertiseAddress)
	}

	if err = setupMetrics(); err != nil {
		return nil, fmt.Errorf("failed to setup metrics: %w", err)
	}

	server.config = config
	server.serfEvents, server.serf, err = setupSerf(config, server.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup serf: %w", err)
	}

	server.raft, server.raftStore, server.raftTransport, err = setupRaft(config, server, server.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup raft: %w", err)
	}

	server.store, err = setupStore(config, server.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	executor := command.NewExecutor(server.raft, config.Raft.Timeout)

	// Node stack
	info := node.Info{
		LocalID: raft.ServerID(config.AdvertiseAddress),
		Version: config.Version,
	}
	server.nodeGRPC = node.NewGRPC(server.raft, info, backup.NewBoltDB(server.store))

	// ACL stack
	server.aclStore = acl.NewBoltStore(server.store)
	server.aclHandler = acl.NewHandler(server.aclStore, server.logger)
	server.aclGRPC = acl.NewGRPC(executor, server.aclStore)

	// Topic stack
	server.topicStore = topic.NewBoltStore(server.store)
	server.topicHandler = topic.NewHandler(server.topicStore, server.logger)
	server.topicGRPC = topic.NewGRPC(executor, server.topicStore)

	// Consumer stack
	server.consumerStore = consumer.NewBoltStore(server.store)
	server.consumerHandler = consumer.NewHandler(server.consumerStore, server.logger)
	server.consumerGRPC = consumer.NewGRPC(executor)

	// Signing stack
	server.signingStore = signing.NewBoltStore(server.store)
	server.signingHandler = signing.NewHandler(server.signingStore, server.logger)
	server.signingGRPC = signing.NewGRPC(executor, server.signingStore)

	// Message stack
	server.messageStore = message.NewBoltStore(server.store)
	server.messageHandler = message.NewHandler(server.messageStore, server.logger)
	server.messageGRPC = message.NewGRPC(
		executor,
		server.messageStore,
		server.consumerStore,
		server.aclStore,
		server.signingStore,
		server.topicStore,
	)

	// Pruning stack
	server.pruner = prune.New(server.topicStore, server.messageStore, server.consumerStore, server.logger)

	return server, nil
}

func setupStore(config Config, logger hclog.Logger) (*bbolt.DB, error) {
	const mode = 0o755
	options := bbolt.DefaultOptions

	snapshot := filepath.Join(config.DataPath, "state_snapshot.db")
	store := filepath.Join(config.DataPath, "state.db")

	_, err := os.Stat(snapshot)
	switch {
	case errors.Is(err, os.ErrNotExist):
		// We don't have a snapshot file waiting to be replaced, so just open the normal store.
		logger.Debug("no restored snapshots found, opening default store")
		return bbolt.Open(store, mode, options)
	case err != nil:
		return nil, err
	}

	logger.Debug("found an state_snapshot.db file, replacing existing store")

	// We've found a snapshot file, remove any existing store, rename the snapshot to match
	// the expected store name and open that.
	if err = os.RemoveAll(store); err != nil {
		return nil, err
	}

	if err = os.Rename(snapshot, store); err != nil {
		return nil, err
	}

	return bbolt.Open(store, mode, options)
}

// ErrReload is the error given when the server has read a snapshot and must be restarted to restore its state.
var ErrReload = errors.New("server has requested to reload from snapshot")

// Start the server. This method blocks until an error occurs or the provided context is cancelled.
func (svr *Server) Start(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return svr.handleSerfEvents(ctx)
	})

	grp.Go(func() error {
		return svr.serveGRPC(ctx)
	})

	grp.Go(func() error {
		return svr.serveMetrics(ctx, svr.config)
	})

	grp.Go(func() error {
		return svr.exportMetrics(ctx)
	})

	grp.Go(func() error {
		return svr.pruner.Prune(ctx, svr.config.PruneInterval)
	})

	grp.Go(func() error {
		<-ctx.Done()
		svr.logger.Info("server shutting down")

		// Shutting down gracefully requires several operations
		return multierror.Append(
			// Leave the serf cluster
			svr.serf.Leave(),
			// Remove ourselves from the raft configuration
			svr.raft.RemoveServer(raft.ServerID(svr.config.AdvertiseAddress), 0, time.Minute).Error(),
			// Shutdown all the raft goroutines
			svr.raft.Shutdown().Error(),
			// Shutdown all the serf goroutines
			svr.serf.Shutdown(),
			// Close and sync the state store
			svr.store.Close(),
		).ErrorOrNil()
	})

	grp.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-svr.restore:
			return ErrReload
		}
	})

	err := grp.Wait()
	switch {
	case errors.Is(err, http.ErrServerClosed), errors.Is(err, grpc.ErrServerStopped), errors.Is(err, raft.ErrRaftShutdown):
		return context.Canceled
	default:
		return err
	}
}

// ErrNoDefaultAdvertiseAddress is the error returned when the server cannot determine an IP address to advertise to
// other servers in the cluster.
var ErrNoDefaultAdvertiseAddress = errors.New("could not find a default advertise address to use")

func getDefaultAdvertiseAddress() (string, error) {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to list interface addresses: %w", err)
	}

	for _, address := range interfaceAddrs {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	return "", ErrNoDefaultAdvertiseAddress
}
