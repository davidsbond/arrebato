package arrebato

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	nodesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/node/service/v1"
)

type (
	// The Node type describes a single node in the cluster.
	Node struct {
		// The name of the node.
		Name string `json:"name"`
		// Whether this node is the leader.
		Leader bool `json:"leader"`
		// The version of the node.
		Version string `json:"version"`
		// Peers known to the node
		Peers []string `json:"peers"`
	}
)

// Nodes returns a slice of all nodes known to the client.
func (c *Client) Nodes(ctx context.Context) ([]Node, error) {
	var nodes []Node
	err := c.cluster.forEach(ctx, func(conn *grpc.ClientConn) error {
		resp, err := nodesvc.NewNodeServiceClient(conn).Describe(ctx, &nodesvc.DescribeRequest{})
		if err != nil {
			return err
		}

		nodes = append(nodes, Node{
			Name:    resp.GetNode().GetName(),
			Leader:  resp.GetNode().GetLeader(),
			Version: resp.GetNode().GetVersion(),
			Peers:   resp.GetNode().GetPeers(),
		})

		return nil
	})

	return nodes, err
}

// ErrNoNode is the error given when querying a node that does not exist.
var ErrNoNode = errors.New("no node")

// Node returns information on a named node. Returns ErrNoNode if the node is not known to the client.
func (c *Client) Node(ctx context.Context, name string) (Node, error) {
	conn, ok := c.cluster.named(name)
	if !ok {
		return Node{}, ErrNoNode
	}

	resp, err := nodesvc.NewNodeServiceClient(conn).Describe(ctx, &nodesvc.DescribeRequest{})
	if err != nil {
		return Node{}, err
	}

	return Node{
		Name:    resp.GetNode().GetName(),
		Leader:  resp.GetNode().GetLeader(),
		Version: resp.GetNode().GetVersion(),
		Peers:   resp.GetNode().GetPeers(),
	}, nil
}
