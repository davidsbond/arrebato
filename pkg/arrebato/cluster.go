package arrebato

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
)

type cluster struct {
	mux         *sync.Mutex
	leaderNode  *grpc.ClientConn
	connections []*grpc.ClientConn
	idx         int
}

func newCluster(ctx context.Context, connections []*grpc.ClientConn) *cluster {
	cl := &cluster{connections: connections, mux: &sync.Mutex{}}
	cl.findLeader(ctx)

	return cl
}

func (c *cluster) any() *grpc.ClientConn {
	c.mux.Lock()
	defer c.mux.Unlock()

	conn := c.connections[c.idx]
	c.idx++
	if c.idx > len(c.connections)-1 {
		c.idx = 0
	}

	return conn
}

func (c *cluster) leader() *grpc.ClientConn {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.leaderNode
}

func (c *cluster) findLeader(ctx context.Context) {
	for _, connection := range c.connections {
		resp, err := nodesvc.NewNodeServiceClient(connection).Describe(ctx, &nodesvc.DescribeRequest{})
		if err != nil {
			continue
		}

		if resp.GetNode().GetLeader() {
			c.mux.Lock()
			c.leaderNode = connection
			c.mux.Unlock()
			return
		}
	}
}

func (c *cluster) Close() error {
	errs := make([]error, 0)
	for _, connection := range c.connections {
		if err := connection.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
