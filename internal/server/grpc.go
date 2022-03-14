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

	"github.com/davidsbond/arrebato/internal/clientinfo"

	// Enable gzip compression from gRPC clients.
	_ "google.golang.org/grpc/encoding/gzip"
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

	svr.consumerGRPC.Register(server)
	svr.messageGRPC.Register(server)
	svr.topicGRPC.Register(server)
	svr.aclGRPC.Register(server)

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
		server.GracefulStop()
		return nil
	})

	return grp.Wait()
}
