package server

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"

	"github.com/davidsbond/arrebato/internal/command"
	nodecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/node/command/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
)

type (
	// The SerfConfig type contains configuration values for serf.
	SerfConfig struct {
		// The Port to use for serf transport.
		Port int
	}
)

func setupSerf(config Config, logger hclog.Logger) (<-chan serf.Event, *serf.Serf, error) {
	memberlistConfig := memberlist.DefaultLANConfig()
	memberlistConfig.BindAddr = config.BindAddress
	memberlistConfig.BindPort = config.Serf.Port
	memberlistConfig.AdvertiseAddr = config.AdvertiseAddress
	memberlistConfig.AdvertisePort = config.Serf.Port
	memberlistConfig.Logger = logger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevelsWithTimestamp: true,
	})

	serfEvents := make(chan serf.Event, 1)
	serfConfig := serf.DefaultConfig()
	serfConfig.Logger = logger.StandardLogger(&hclog.StandardLoggerOptions{
		InferLevelsWithTimestamp: true,
	})

	serfConfig.EventCh = serfEvents
	serfConfig.NodeName = config.AdvertiseAddress
	serfConfig.SnapshotPath = filepath.Join(config.DataPath, "serf.db")
	serfConfig.RejoinAfterLeave = true

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

		peer := fmt.Sprint(member.Name, ":", svr.config.Raft.Port)
		err := svr.raft.AddVoter(serverID, raft.ServerAddress(peer), 0, svr.config.Raft.Timeout).Error()
		switch {
		case errors.Is(err, raft.ErrLeadershipLost):
			return nil
		case err != nil:
			return fmt.Errorf("failed to add server: %w", err)
		}

		cmd := command.New(&nodecmd.AddNode{
			Node: &node.Node{
				Name: member.Name,
			},
		})

		if err = svr.executor.Execute(ctx, cmd); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
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

		svr.logger.Info("removing server", "name", member.Name)
		err := svr.raft.RemoveServer(raft.ServerID(member.Name), 0, svr.config.Raft.Timeout).Error()
		switch {
		case errors.Is(err, raft.ErrLeadershipLost):
			return nil
		case err != nil:
			return fmt.Errorf("failed to remove existing server: %w", err)
		}

		cmd := command.New(&nodecmd.RemoveNode{
			Name: member.Name,
		})

		if err = svr.executor.Execute(ctx, cmd); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return nil
}
