package arrebato

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
)

type cluster struct {
	mux        *sync.RWMutex
	leaderNode *grpc.ClientConn
	nodes      map[string]*grpc.ClientConn
}

func newCluster(ctx context.Context, connections []*grpc.ClientConn) (*cluster, error) {
	nodes := make(map[string]*grpc.ClientConn)
	for _, connection := range connections {
		resp, err := nodesvc.NewNodeServiceClient(connection).Describe(ctx, &nodesvc.DescribeRequest{})
		if err != nil {
			return nil, err
		}

		nodes[resp.GetNode().GetName()] = connection
	}

	cl := &cluster{nodes: nodes, mux: &sync.RWMutex{}}
	cl.findLeader(ctx)

	return cl, nil
}

func (c *cluster) any() *grpc.ClientConn {
	c.mux.RLock()
	defer c.mux.RUnlock()

	// Map ordering is non-deterministic, we can just pick the first one we get
	// in the range.
	for _, conn := range c.nodes {
		return conn
	}

	return nil
}

func (c *cluster) leader() *grpc.ClientConn {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.leaderNode
}

func (c *cluster) named(name string) (*grpc.ClientConn, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	conn, ok := c.nodes[name]
	return conn, ok
}

func (c *cluster) topicOwner(ctx context.Context, topic string) (*grpc.ClientConn, error) {
	var selected *grpc.ClientConn
	err := c.forEach(ctx, func(conn *grpc.ClientConn) error {
		if selected != nil {
			return nil
		}

		n, err := nodesvc.NewNodeServiceClient(conn).Describe(ctx, &nodesvc.DescribeRequest{})
		if err != nil {
			return err
		}

		for _, tp := range n.GetNode().GetTopics() {
			if tp == topic {
				selected = conn
				return nil
			}
		}

		return nil
	})

	switch {
	case err != nil:
		return nil, err
	case selected == nil:
		return nil, ErrNoTopic
	default:
		return selected, nil
	}
}

func (c *cluster) findLeader(ctx context.Context) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, connection := range c.nodes {
		resp, err := nodesvc.NewNodeServiceClient(connection).Describe(ctx, &nodesvc.DescribeRequest{})
		if err != nil {
			continue
		}

		if resp.GetNode().GetLeader() {
			c.leaderNode = connection
			return
		}
	}
}

func (c *cluster) forEach(ctx context.Context, fn func(conn *grpc.ClientConn) error) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	for _, connection := range c.nodes {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := fn(connection); err != nil {
			return err
		}
	}

	return nil
}

func (c *cluster) Close() error {
	errs := make([]error, 0)
	for _, connection := range c.nodes {
		if err := connection.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}
