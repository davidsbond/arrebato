package server

import (
	"compress/gzip"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	transport "github.com/Jille/raft-grpc-transport"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/command"
	aclcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/command/v1"
	consumercmd "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/command/v1"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	signingcmd "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/command/v1"
	topiccmd "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/command/v1"
)

type (
	// The RaftConfig type describes configuration values for the raft consensus algorithm used to maintain state
	// across the cluster.
	RaftConfig struct {
		// Timeout is the timeout to use for raft communications.
		Timeout time.Duration

		// MaxSnapshots is the maximum number of raft snapshots to keep.
		MaxSnapshots int

		// NonVoter determines if the server is added to the cluster as a replica that can
		// never gain leadership.
		NonVoter bool
	}
)

const (
	raftLogCacheSize = 512
)

func setupRaft(config Config, fsm raft.FSM, logger hclog.Logger) (*raft.Raft, *raftboltdb.BoltStore, *transport.Manager, error) {
	if err := os.MkdirAll(config.DataPath, 0o750); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create raft data dir: %w", err)
	}

	raftAddress := fmt.Sprint(config.BindAddress, ":", config.GRPC.Port)
	options := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithCodes(codes.Unavailable),
				grpc_retry.WithMax(100),
			),
			grpc_prometheus.UnaryClientInterceptor,
			clientinfo.UnaryClientInterceptor(config.AdvertiseAddress),
		),
		grpc.WithChainStreamInterceptor(
			grpc_retry.StreamClientInterceptor(),
			grpc_prometheus.StreamClientInterceptor,
			clientinfo.StreamClientInterceptor(config.AdvertiseAddress),
		),
	}

	if config.TLS.enabled() {
		// We generate a TLS client configuration here to enable TLS for raft transport.
		tlsConfig, err := config.TLS.tlsConfig()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create TLS config: %w", err)
		}

		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	raftTransport := transport.New(raft.ServerAddress(raftAddress), options)
	raftConfig := raft.DefaultConfig()
	raftConfig.Logger = logger
	raftConfig.LocalID = raft.ServerID(config.AdvertiseAddress)

	dataPath := filepath.Join(config.DataPath, "raft.db")
	db, err := raftboltdb.NewBoltStore(dataPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create database: %w", err)
	}

	snapshots, err := raft.NewFileSnapshotStoreWithLogger(config.DataPath, config.Raft.MaxSnapshots, logger)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create snapshot store: %w", err)
	}

	// We wrap the log store with a cache to increase performance accessing recent logs.
	logs, err := raft.NewLogCache(raftLogCacheSize, db)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create log cache: %w", err)
	}

	r, err := raft.NewRaft(raftConfig, fsm, logs, db, snapshots, raftTransport.Transport())
	if err != nil {
		return nil, nil, nil, err
	}

	suffrage := raft.Voter
	if config.Raft.NonVoter {
		suffrage = raft.Nonvoter
	}

	bootstrap := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: suffrage,
				ID:       raftConfig.LocalID,
				Address:  raft.ServerAddress(raftAddress),
			},
		},
	}

	for _, peer := range config.Peers {
		if peer == config.AdvertiseAddress {
			continue
		}

		peerAddress := fmt.Sprint(peer, ":", config.GRPC.Port)
		bootstrap.Servers = append(bootstrap.Servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(peer),
			Address:  raft.ServerAddress(peerAddress),
		})
	}

	if ok, _ := raft.HasExistingState(db, db, snapshots); ok {
		logger.Debug("found existing raft state")
		return r, db, raftTransport, nil
	}

	logger.Debug("bootstrapping cluster")
	err = r.BootstrapCluster(bootstrap).Error()
	switch {
	case errors.Is(err, raft.ErrCantBootstrap):
		// We get this error if a raft state already exists in the data path, this usually means that this node
		break
	case err != nil:
		return nil, nil, nil, fmt.Errorf("failed to bootstrap: %w", err)
	}

	return r, db, raftTransport, nil
}

// Snapshot returns a raft.FSMSnapshot implementation that backs up the current state of the applied raft log. It
// always returns a nil error.
func (svr *Server) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{store: svr.store}, nil
}

// Restore replaces the current state of the applied raft log with the contents of the io.ReadCloser implementation. This
// is done by initialising the snapshot in a temporary file and triggering the server to restart. The server will detect
// the restore file, rename it and use it from then on.
func (svr *Server) Restore(snapshot io.ReadCloser) error {
	restorePath := filepath.Join(svr.config.DataPath, "state_snapshot.db")

	file, err := os.Create(restorePath)
	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(snapshot)
	if err != nil {
		return fmt.Errorf("failed to read gzipped snapshot: %w", err)
	}

	for {
		_, err = io.CopyN(file, reader, 1024)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to copy snapshot to temporary file: %w", err)
		}
	}

	if err = reader.Close(); err != nil {
		return fmt.Errorf("failed to close gzip reader: %w", err)
	}

	if err = snapshot.Close(); err != nil {
		return fmt.Errorf("failed to close snapshot: %w", err)
	}

	svr.logger.Info("server has restored a snapshot, restarting")
	svr.restore <- struct{}{}
	return nil
}

// Apply unmarshals the contents of the raft.Log, expecting a command.Command that can be handled by a command
// handler.
func (svr *Server) Apply(log *raft.Log) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), svr.config.Raft.Timeout)
	defer cancel()

	lastAppliedIndex, err := svr.lastAppliedIndex()
	if err != nil {
		return fmt.Errorf("failed to get last applied index: %w", err)
	}

	// When the server restarts, the log will be replayed, we don't want to duplicate all the messages/topics in
	// the state so if the log index is less than the last known index sent to the FSM then we do nothing. We can't
	// fully rely on the raft mechanism to know exactly the last log index that the FSM successfully handled, so we
	// also track that manually.
	if log.Index <= lastAppliedIndex && lastAppliedIndex != 0 {
		return nil
	}

	cmd, err := command.Unmarshal(log.Data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal command: %w", err)
	}

	switch payload := cmd.Payload().(type) {
	case *topiccmd.CreateTopic:
		err = svr.topicHandler.Create(ctx, payload)
	case *topiccmd.DeleteTopic:
		err = svr.topicHandler.Delete(ctx, payload)
	case *messagecmd.CreateMessage:
		err = svr.messageHandler.Create(ctx, payload)
	case *consumercmd.SetTopicIndex:
		err = svr.consumerHandler.SetTopicIndex(ctx, payload)
	case *aclcmd.SetACL:
		err = svr.aclHandler.Set(ctx, payload)
	case *signingcmd.CreatePublicKey:
		err = svr.signingHandler.Create(ctx, payload)
	case *nodecmd.AddNode:
		err = svr.nodeHandler.Add(ctx, payload)
	case *nodecmd.RemoveNode:
		err = svr.nodeHandler.Remove(ctx, payload)
	case *nodecmd.AssignTopic:
		err = svr.nodeHandler.AssignTopic(ctx, payload)
	default:
		break
	}

	if err != nil {
		return err
	}

	return svr.setLastAppliedIndex(log.Index)
}

const (
	raftKey             = "raft"
	lastAppliedIndexKey = "last_applied_index"
)

func (svr *Server) setLastAppliedIndex(index uint64) error {
	return svr.store.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(raftKey))
		if err != nil {
			return fmt.Errorf("failed to open raft bucket: %w", err)
		}

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, index)

		return bucket.Put([]byte(lastAppliedIndexKey), value)
	})
}

func (svr *Server) lastAppliedIndex() (uint64, error) {
	var index uint64
	err := svr.store.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(raftKey))
		if bucket == nil {
			return nil
		}

		value := bucket.Get([]byte(lastAppliedIndexKey))
		if value == nil {
			return nil
		}

		index = binary.BigEndian.Uint64(value)
		return nil
	})

	return index, err
}

// IsLeader returns true if this server instance is the cluster leader.
func (svr *Server) IsLeader() bool {
	return svr.raft.State() == raft.Leader
}
