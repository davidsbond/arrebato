package node_test

import (
	"testing"

	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/node"
	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_Describe(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		State    raft.RaftState
		Config   raft.Configuration
		Expected *nodesvc.DescribeResponse
	}{
		{
			Name:  "It should return node details and known peers",
			State: raft.Leader,
			Config: raft.Configuration{
				Servers: []raft.Server{
					{
						ID: raft.ServerID("test-server"),
					},
					{
						ID: raft.ServerID("test-server-2"),
					},
				},
			},
			Expected: &nodesvc.DescribeResponse{
				Node: &nodepb.Node{
					Name:    "test-server",
					Leader:  true,
					Peers:   []string{"test-server-2"},
					Version: "v1.0.0",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			r := &MockRaft{state: tc.State, config: tc.Config}
			localID := raft.ServerID("test-server")
			actual, err := node.NewGRPC(r, localID, "v1.0.0").Describe(ctx, &nodesvc.DescribeRequest{})
			assert.NoError(t, err)
			assert.True(t, proto.Equal(tc.Expected, actual))
		})
	}
}
