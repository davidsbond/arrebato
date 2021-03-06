package server

import (
	"context"
	"fmt"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/davidsbond/arrebato/internal/clientinfo"

	// Enable gzip compression from gRPC clients.
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type (
	// The GRPCConfig type describes configuration values for the Server's gRPC endpoints.
	GRPCConfig struct {
		// The Port to use for gRPC transport.
		Port int
	}
)

func (svr *Server) serveGRPC(ctx context.Context) error {
	address := fmt.Sprint(svr.config.BindAddress, ":", svr.config.GRPC.Port)

	grpc_prometheus.EnableHandlingTimeHistogram()

	var options []grpc.ServerOption
	infoExtractor := clientinfo.MetadataExtractor

	if svr.config.TLS.enabled() {
		config, err := svr.config.TLS.tlsConfig()
		if err != nil {
			return fmt.Errorf("failed to create tls config: %w", err)
		}

		options = append(options, grpc.Creds(credentials.NewTLS(config)))
		infoExtractor = clientinfo.TLSExtractor
	} else {
		svr.logger.Warn("no TLS credentials were provided, gRPC transport will be insecure")
	}

	options = append(options,
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			clientinfo.UnaryServerInterceptor(infoExtractor),
		),
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			clientinfo.StreamServerInterceptor(infoExtractor),
		),
	)

	server := grpc.NewServer(options...)
	healthServer := health.NewServer()

	svr.consumerGRPC.Register(server, healthServer)
	svr.messageGRPC.Register(server, healthServer)
	svr.topicGRPC.Register(server, healthServer)
	svr.aclGRPC.Register(server, healthServer)
	svr.nodeGRPC.Register(server, healthServer)
	svr.signingGRPC.Register(server, healthServer)
	svr.raftTransport.Register(server)

	// Part of the gRPC health checking protocol states that the server should use an empty string as the key for
	// the server's overall health status.
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("failed to bind grpc listener: %w", err)
		}

		svr.logger.Named("grpc").Info("starting listener", "address", address)

		return server.Serve(listener)
	})

	grp.Go(func() error {
		<-ctx.Done()
		healthServer.Shutdown()
		server.GracefulStop()
		return nil
	})

	return grp.Wait()
}
