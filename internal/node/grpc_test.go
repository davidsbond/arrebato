package node_test

import (
	"io"
	"testing"

	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/node"
	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	"github.com/davidsbond/arrebato/internal/testutil"
)

func TestGRPC_Describe(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		State        raft.RaftState
		Nodes        []*nodepb.Node
		Error        error
		ExpectedCode codes.Code
		Expected     *nodesvc.DescribeResponse
	}{
		{
			Name:  "It should return node details and known peers",
			State: raft.Leader,
			Nodes: []*nodepb.Node{
				{
					Name:   "test-server",
					Topics: []string{"test-topic"},
				},
				{
					Name: "test-server-2",
				},
			},
			Expected: &nodesvc.DescribeResponse{
				Node: &nodepb.Node{
					Name:    "test-server",
					Leader:  true,
					Peers:   []string{"test-server-2"},
					Version: "v1.0.0",
					Topics:  []string{"test-topic"},
				},
			},
		},
		{
			Name:         "It should return errors to the caller",
			Error:        io.EOF,
			ExpectedCode: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := testutil.Context(t)
			r := &MockRaft{state: tc.State}
			lister := &MockLister{nodes: tc.Nodes, err: tc.Error}
			info := node.Info{
				Name:    "test-server",
				Version: "v1.0.0",
			}

			actual, err := node.NewGRPC(r, info, nil, lister).Describe(ctx, &nodesvc.DescribeRequest{})
			assert.Equal(t, tc.ExpectedCode, status.Code(err))

			if tc.Error != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(tc.Expected, actual))
		})
	}
}
