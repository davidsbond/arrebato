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

	"github.com/davidsbond/arrebato/internal/clientinfo"

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
		// Addresses of the running servers, multiple addresses are expected here so that the client can correctly
		// route write operations to the leader.
		Addresses []string

		// Configuration for connecting to the server via TLS. When using TLS, the server will expect clients to be
		// issued a SPIFFE ID for identification.
		TLS *tls.Config

		// The identifier for the client, this is only required when running the server in an insecure mode, when using
		// TLS, it is expected that the client certificate will contain a SPIFFE ID that the client will use to
		// identify itself.
		ClientID string
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
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
		grpc.WithChainUnaryInterceptor(
			clientinfo.UnaryClientInterceptor(config.ClientID),
		),
		grpc.WithChainStreamInterceptor(
			clientinfo.StreamClientInterceptor(config.ClientID),
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

	cl, err := newCluster(ctx, connections)
	if err != nil {
		return nil, err
	}

	return &Client{
		cluster: cl,
		config:  config,
	}, nil
}

// Close the connection to the cluster.
func (c *Client) Close() error {
	return c.cluster.Close()
}
