package e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/davidsbond/arrebato/internal/server"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type Suite struct {
	suite.Suite

	cancel context.CancelFunc
	server *server.Server
	client *arrebato.Client
	config server.Config
}

func (s *Suite) SetupSuite() {
	tmp, err := os.MkdirTemp(os.TempDir(), "arrebato-e2e")
	require.NoError(s.T(), err)

	s.config = server.Config{
		LogLevel:      10,
		BindAddress:   "0.0.0.0",
		DataPath:      tmp,
		PruneInterval: time.Minute,
		Raft: server.RaftConfig{
			Timeout:      time.Minute,
			MaxSnapshots: 3,
		},
		Serf: server.SerfConfig{
			Port: 5001,
		},
		GRPC: server.GRPCConfig{
			Port: 5000,
		},
		Metrics: server.MetricConfig{
			Port: 5002,
		},
	}

	s.server, err = server.New(s.config)
	require.NoError(s.T(), err)

	s.T().Log("starting test server")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		_ = s.server.Start(ctx)
	}()

	s.T().Log("waiting for test server to become leader")
	require.Eventually(s.T(), func() bool {
		return s.server.IsLeader()
	}, time.Minute, time.Second)

	s.cancel = cancel
	s.client, err = arrebato.Dial(ctx, arrebato.Config{
		Addresses: []string{
			":5000",
		},
		ClientID: "test-client",
	})

	require.NoError(s.T(), err)
}

func (s *Suite) TearDownSuite() {
	s.T().Log("stopping test server")
	s.cancel()

	assert.NoError(s.T(), s.client.Close())
	require.Eventually(s.T(), func() bool {
		return !s.server.IsLeader()
	}, time.Minute, time.Second)
}

func (s *Suite) TestSuite() {
	suites := []suite.TestingSuite{
		&TopicSuite{
			client: s.client,
		},
		&MessageSuite{
			client: s.client,
		},
		&SigningSuite{
			client: s.client,
		},
		&ACLSuite{
			client: s.client,
		},
		// Keep the SnapshotSuite last as it requires stopping and restarting the
		// server. Just makes it less complicated for all the other test suites.
		&SnapshotSuite{
			server: s.server,
			config: s.config,
		},
	}

	for _, su := range suites {
		suite.Run(s.T(), su)
	}
}

func Test_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e tests skipped in short mode")
		return
	}

	s := new(Suite)
	suite.Run(t, s)
}
