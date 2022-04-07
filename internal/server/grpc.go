package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// Enable gzip compression from gRPC clients.
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/davidsbond/arrebato/internal/clientinfo"
)

type (
	// The GRPCConfig type describes configuration values for the Server's gRPC endpoints.
	GRPCConfig struct {
		// The Port to use for gRPC transport.
		Port int

		// Location of the TLS certificate file to use for transport credentials.
		TLSCertFile string

		// Location of the TLS Key file to use for transport credentials.
		TLSKeyFile string

		// The certificate of the CA that signs client certificates.
		TLSCAFile string
	}
)

func (c GRPCConfig) tlsEnabled() bool {
	return c.TLSCertFile != "" && c.TLSKeyFile != ""
}

func (svr *Server) serveGRPC(ctx context.Context) error {
	address := fmt.Sprint(svr.config.BindAddress, ":", svr.config.GRPC.Port)

	grpc_prometheus.EnableHandlingTimeHistogram()

	var options []grpc.ServerOption
	infoExtractor := clientinfo.MetadataExtractor

	if svr.config.GRPC.tlsEnabled() {
		ca, err := ioutil.ReadFile(svr.config.GRPC.TLSCAFile)
		if err != nil {
			return fmt.Errorf("failed to read CA file: %w", err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(ca)

		cert, err := tls.LoadX509KeyPair(svr.config.GRPC.TLSCertFile, svr.config.GRPC.TLSKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS files: %w", err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
			MinVersion:   tls.VersionTLS13,
		}

		options = append(options, grpc.Creds(credentials.NewTLS(config)))
		infoExtractor = clientinfo.TLSExtractor
	} else {
		svr.logger.Warn("no TLS credentials were provided, gRPC transport will be insecure")
	}

	options = append(options,
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			clientinfo.UnaryServerInterceptor(infoExtractor),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
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
	svr.raftGRPC.Register(server, healthServer)

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
