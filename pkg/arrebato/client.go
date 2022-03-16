// Package arrebato provides the arrebato client implementation used to interact with a cluster. This includes topic management,
// message consumption and production.
package arrebato

import (
	"context"
	"crypto/tls"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	// Enable gzip compression for gRPC clients.
	_ "google.golang.org/grpc/encoding/gzip"
)

type (
	// The Client type is used to interact with an arrebato cluster.
	Client struct {
		cluster *cluster
		config  Config
	}

	// The Config type describes configuration values used by a Client.
	Config struct {
		Addresses         []string
		TLS               *tls.Config
		ClientID          string
		MessageSigningKey []byte
	}
)

// DefaultConfig returns a Config instance with sane values for a Client's connection.
func DefaultConfig(addrs []string) Config {
	return Config{
		Addresses: addrs,
		ClientID:  uuid.New().String(),
	}
}

// Dial an arrebato cluster, returning a Client that can be used to perform requests against it.
func Dial(ctx context.Context, config Config) (*Client, error) {
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
		grpc.WithChainUnaryInterceptor(
			unaryClientInterceptor(config.ClientID),
		),
		grpc.WithChainStreamInterceptor(
			streamClientInterceptor(config.ClientID),
		),
	}

	if config.TLS != nil {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(config.TLS)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	connections := make([]*grpc.ClientConn, len(config.Addresses))
	for i, address := range config.Addresses {
		conn, err := grpc.DialContext(ctx, address, options...)
		if err != nil {
			return nil, err
		}

		connections[i] = conn
	}

	cl := newCluster(ctx, connections)

	return &Client{
		cluster: cl,
		config:  config,
	}, nil
}

// Close the connection to the cluster.
func (c *Client) Close() error {
	return c.cluster.Close()
}

func unaryClientInterceptor(clientID string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Client-ID", clientID)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func streamClientInterceptor(clientID string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Client-ID", clientID)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
