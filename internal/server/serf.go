package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

type (
	// The SerfConfig type contains configuration values for serf.
	SerfConfig struct {
		// The Port to use for serf transport.
		Port int

		// The location of the file whose contents should contain the primary encryption key for
		// gossip messages. The file contents should be either 16, 24, or 32 bytes to select AES-128,
		// AES-192, or AES-256.
		EncryptionKeyFile string
	}
)

const (
	voterKey            = "voter"
	grpcPortKey         = "grpc_port"
	roleKey             = "role"
	advertiseAddressKey = "advertise_address"
)

func setupSerf(config Config, logger hclog.Logger) (<-chan serf.Event, *serf.Serf, error) {
	var secretKey []byte
	if config.Serf.EncryptionKeyFile != "" {
		data, err := os.ReadFile(config.Serf.EncryptionKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read encryption key file: %w", err)
		}

		secretKey = data
	}

	if len(secretKey) == 0 {
		logger.Warn("serf will not use an encryption key, gossip will be insecure")
	}

	memberlistConfig := memberlist.DefaultLANConfig()
	memberlistConfig.BindAddr = config.BindAddress
	memberlistConfig.BindPort = config.Serf.Port
	memberlistConfig.AdvertiseAddr = config.AdvertiseAddress
	memberlistConfig.AdvertisePort = config.Serf.Port
	memberlistConfig.Logger = logger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevels:              true,
		InferLevelsWithTimestamp: true,
	})
	if len(secretKey) > 0 {
		memberlistConfig.SecretKey = secretKey
		memberlistConfig.GossipVerifyIncoming = true
		memberlistConfig.GossipVerifyOutgoing = true
	}

	serfEvents := make(chan serf.Event, 1)
	serfConfig := serf.DefaultConfig()
	serfConfig.Logger = logger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevels:              true,
		InferLevelsWithTimestamp: true,
	})

	// There are multiple serf files that need persisting, we place them under their own
	// directory in the data path. We make sure that exists first.
	if err := os.MkdirAll(filepath.Join(config.DataPath, "serf"), 0o750); err != nil {
		return nil, nil, fmt.Errorf("failed to create serf directory: %w", err)
	}

	serfConfig.EventCh = serfEvents
	serfConfig.NodeName = config.AdvertiseAddress
	serfConfig.SnapshotPath = filepath.Join(config.DataPath, "serf", "serf.db")
	serfConfig.KeyringFile = filepath.Join(config.DataPath, "serf", "serf.keyring")
	serfConfig.RejoinAfterLeave = true
	serfConfig.Tags = map[string]string{
		voterKey:            strconv.FormatBool(!config.Raft.NonVoter),
		grpcPortKey:         strconv.Itoa(config.GRPC.Port),
		advertiseAddressKey: config.AdvertiseAddress,
		roleKey:             "server",
	}

	s, err := serf.Create(serfConfig)
	if err != nil {
		return nil, nil, err
	}

	peers := make([]string, 0)
	for _, peer := range config.Peers {
		if peer == config.AdvertiseAddress {
			continue
		}

		peers = append(peers, peer)
	}

	if len(peers) > 0 {
		if _, err = s.Join(config.Peers, false); err != nil {
			return nil, nil, fmt.Errorf("failed to join: %w", err)
		}
	}

	return serfEvents, s, nil
}

func (svr *Server) handleSerfEvents(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-svr.serfEvents:
			if err := svr.handleSerfEvent(ctx, event); err != nil {
				return fmt.Errorf("failed to handle serf event: %w", err)
			}
		}
	}
}

func (svr *Server) handleSerfEvent(ctx context.Context, event serf.Event) error {
	switch payload := event.(type) {
	case serf.MemberEvent:
		return svr.handleSerfMemberEvent(ctx, payload)
	default:
		return nil
	}
}

func (svr *Server) handleSerfMemberEvent(ctx context.Context, event serf.MemberEvent) error {
	switch event.EventType() {
	case serf.EventMemberJoin:
		return svr.handleSerfEventMemberJoin(ctx, event)
	case serf.EventMemberLeave, serf.EventMemberFailed, serf.EventMemberReap:
		return svr.handleSerfEventMemberLeave(ctx, event)
	default:
		return nil
	}
}

func (svr *Server) handleSerfEventMemberJoin(ctx context.Context, event serf.MemberEvent) error {
	if err := svr.raft.VerifyLeader().Error(); err != nil {
		return nil
	}

	for _, member := range event.Members {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !isServer(member.Tags) {
			continue
		}

		future := svr.raft.GetConfiguration()
		if future.Error() != nil {
			return fmt.Errorf("failed to get raft configuration: %w", future.Error())
		}

		serverID := raft.ServerID(member.Name)
		for _, server := range future.Configuration().Servers {
			if serverID != server.ID {
				continue
			}

			err := svr.raft.RemoveServer(serverID, 0, svr.config.Raft.Timeout).Error()
			switch {
			case errors.Is(err, raft.ErrLeadershipLost):
				return nil
			case err != nil:
				return fmt.Errorf("failed to remove existing server: %w", err)
			}
		}

		var err error
		voter := isVoter(member.Tags)
		peer := net.JoinHostPort(member.Tags[advertiseAddressKey], member.Tags[grpcPortKey])

		if voter {
			err = svr.raft.AddVoter(serverID, raft.ServerAddress(peer), 0, svr.config.Raft.Timeout).Error()
		} else {
			err = svr.raft.AddNonvoter(serverID, raft.ServerAddress(peer), 0, svr.config.Raft.Timeout).Error()
		}

		switch {
		case errors.Is(err, raft.ErrLeadershipLost):
			return nil
		case err != nil:
			return fmt.Errorf("failed to add server: %w", err)
		}
	}

	return nil
}

func (svr *Server) handleSerfEventMemberLeave(ctx context.Context, event serf.MemberEvent) error {
	if err := svr.raft.VerifyLeader().Error(); err != nil {
		return nil
	}

	for _, member := range event.Members {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !isServer(member.Tags) {
			continue
		}

		svr.logger.Info("removing server", "name", member.Name)
		err := svr.raft.RemoveServer(raft.ServerID(member.Name), 0, svr.config.Raft.Timeout).Error()
		switch {
		case errors.Is(err, raft.ErrLeadershipLost):
			return nil
		case err != nil:
			return fmt.Errorf("failed to remove existing server: %w", err)
		}
	}

	return nil
}

func isVoter(m map[string]string) bool {
	return m[voterKey] == "true"
}

func isServer(m map[string]string) bool {
	return m[roleKey] == "server"
}
