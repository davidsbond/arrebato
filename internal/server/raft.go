package server

import (
	"compress/gzip"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"go.etcd.io/bbolt"

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
		// Port is the port to use for raft transport.
		Port int

		// Timeout is the timeout to use for raft communications.
		Timeout time.Duration

		// MaxPool is the maximum number of connections in the TCP pool.
		MaxPool int

		// MaxSnapshots is the maximum number of raft snapshots to keep.
		MaxSnapshots int
	}
)

const (
	raftLogCacheSize = 512
)

func setupRaft(config Config, fsm raft.FSM, logger hclog.Logger) (*raft.Raft, *raftboltdb.BoltStore, error) {
	if err := os.MkdirAll(config.DataPath, 0o750); err != nil {
		return nil, nil, fmt.Errorf("failed to create raft data dir: %w", err)
	}

	raftAddress := fmt.Sprint(config.BindAddress, ":", config.Raft.Port)
	addrs, err := net.LookupIP(config.AdvertiseAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to look up ip for advertise address: %w", err)
	}

	var advertiseAddress net.TCPAddr
	for _, addr := range addrs {
		if addr.IsLoopback() || addr.To4() == nil {
			continue
		}

		advertiseAddress.IP = addr
		advertiseAddress.Port = config.Raft.Port
		break
	}

	transport, err := raft.NewTCPTransportWithLogger(
		raftAddress,
		&advertiseAddress,
		config.Raft.MaxPool,
		config.Raft.Timeout,
		logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create TCP transport: %w", err)
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.Logger = logger
	raftConfig.LocalID = raft.ServerID(config.AdvertiseAddress)

	dataPath := filepath.Join(config.DataPath, "raft.db")
	db, err := raftboltdb.NewBoltStore(dataPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database: %w", err)
	}

	snapshots, err := raft.NewFileSnapshotStoreWithLogger(config.DataPath, config.Raft.MaxSnapshots, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create snapshot store: %w", err)
	}

	// We wrap the log store with a cache to increase performance accessing recent logs.
	logs, err := raft.NewLogCache(raftLogCacheSize, db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log cache: %w", err)
	}

	r, err := raft.NewRaft(raftConfig, fsm, logs, db, snapshots, transport)
	if err != nil {
		return nil, nil, err
	}

	bootstrap := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raftConfig.LocalID,
				Address:  raft.ServerAddress(raftAddress),
			},
		},
	}

	for _, peer := range config.Peers {
		if peer == config.AdvertiseAddress {
			continue
		}

		peerAddress := fmt.Sprint(peer, ":", config.Raft.Port)
		bootstrap.Servers = append(bootstrap.Servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(peer),
			Address:  raft.ServerAddress(peerAddress),
		})
	}

	if ok, _ := raft.HasExistingState(db, db, snapshots); ok {
		logger.Debug("found existing raft state")
		return r, db, nil
	}

	logger.Debug("bootstrapping cluster")
	err = r.BootstrapCluster(bootstrap).Error()
	switch {
	case errors.Is(err, raft.ErrCantBootstrap):
		// We get this error if a raft state already exists in the data path, this usually means that this node
		break
	case err != nil:
		return nil, nil, fmt.Errorf("failed to bootstrap: %w", err)
	}

	return r, db, nil
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
	restorePath := filepath.Join(svr.config.DataPath, "arrebato_restore.db")

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
	if log.Index < lastAppliedIndex && lastAppliedIndex != 0 {
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
