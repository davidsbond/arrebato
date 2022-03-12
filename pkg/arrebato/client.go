// Package arrebato provides the arrebato client implementation used to interact with a cluster. This includes topic management,
// message consumption and production.
package arrebato

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	// Enable gzip compression for gRPC clients.
	_ "google.golang.org/grpc/encoding/gzip"

	consumersvc "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/service/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	topicsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/service/v1"
)

type (
	// The Client type is used to interact with an arrebato cluster.
	Client struct {
		conn      *grpc.ClientConn
		config    Config
		topics    topicsvc.TopicServiceClient
		messages  messagesvc.MessageServiceClient
		consumers consumersvc.ConsumerServiceClient
	}

	// The Config type describes configuration values used by a Client.
	Config struct {
		Address string
		TLS     *tls.Config
	}
)

// DefaultConfig returns a Config instance with sane values for a Client's connection.
func DefaultConfig(addr string) Config {
	return Config{
		Address: addr,
	}
}

// Dial an arrebato cluster, returning a Client that can be used to perform requests against it.
func Dial(ctx context.Context, config Config) (*Client, error) {
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
	}

	if config.TLS != nil {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(config.TLS)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.DialContext(ctx, config.Address, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      conn,
		topics:    topicsvc.NewTopicServiceClient(conn),
		messages:  messagesvc.NewMessageServiceClient(conn),
		consumers: consumersvc.NewConsumerServiceClient(conn),
		config:    config,
	}, nil
}

// Close the connection to the cluster.
func (c *Client) Close() error {
	return c.conn.Close()
}
