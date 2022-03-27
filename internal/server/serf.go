package server

import (
	"context"
	"fmt"

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
			if err := svr.raft.VerifyLeader().Error(); err != nil {
				continue
			}

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
	for _, member := range event.Members {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		peer := fmt.Sprint(member.Name, ":", svr.config.Raft.Port)
		if err := svr.raft.AddVoter(raft.ServerID(member.Name), raft.ServerAddress(peer), 0, svr.config.Raft.Timeout).Error(); err != nil {
			return fmt.Errorf("failed to add voter %s: %w", peer, err)
		}
	}

	return nil
}

func (svr *Server) handleSerfEventMemberLeave(ctx context.Context, event serf.MemberEvent) error {
	for _, member := range event.Members {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		svr.logger.Info("removing server", "name", member.Name)
		if err := svr.raft.RemoveServer(raft.ServerID(member.Name), 0, svr.config.Raft.Timeout).Error(); err != nil {
			return fmt.Errorf("failed to remove voter %s: %w", member.Name, err)
		}
	}

	return nil
}
