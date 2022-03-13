package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/davidsbond/arrebato/internal/server"
)

var version string

func main() {
	config := server.Config{
		Version: version,
	}

	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Starts a server node",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			undo, err := maxprocs.Set()
			if err != nil {
				return fmt.Errorf("failed to automatically set GOMAXPROCS: %w", err)
			}
			defer undo()

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			svr, err := server.New(config)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			return svr.Start(ctx)
		},
	}

	flags := cmd.PersistentFlags()

	// Server configuration
	flags.StringSliceVar(&config.Peers, "peers", nil, "Existing peers to connect to")
	flags.StringVar(&config.BindAddress, "bind-address", "0.0.0.0", "The network interface to bind to for transport")
	flags.StringVar(&config.AdvertiseAddress, "advertise-address", "", "The IP address to advertise to other nodes")
	flags.StringVar(&config.DataPath, "data-path", "/var/lib/arrebato", "The location to store state on-disk")
	flags.IntVar(&config.LogLevel, "log-level", 3, "The level of logs to produce output for, lower numbers produce more logs")
	flags.DurationVar(&config.PruneInterval, "prune-interval", time.Minute, "The amount of time between message pruning runs")

	// Raft configuration
	flags.IntVar(&config.Raft.Port, "raft-port", 5000, "The port to use for raft transport")
	flags.DurationVar(&config.Raft.Timeout, "raft-timeout", time.Minute, "The timeout to use for raft transport")
	flags.IntVar(&config.Raft.MaxSnapshots, "raft-max-snapshots", 3, "The maximum number of raft snapshots to store")
	flags.IntVar(&config.Raft.MaxPool, "raft-max-pool", 3, "The maximum number of connections in the raft connection pool")

	// Serf configuration
	flags.IntVar(&config.Serf.Port, "serf-port", 5001, "The port to use for serf transport")

	// gRPC configuration
	flags.IntVar(&config.GRPC.Port, "grpc-port", 5002, "The port to use for gRPC transport")
	flags.StringVar(&config.GRPC.TLSCertFile, "grpc-tls-cert", "", "The location of the TLS cert file for the server to use")
	flags.StringVar(&config.GRPC.TLSKeyFile, "grpc-tls-key", "", "The location of the TLS key file for the server to use")
	flags.StringVar(&config.GRPC.TLSCAFile, "grpc-tls-ca", "", "The location of the CA certificate that signs client certificates")

	// Prometheus configuration
	flags.IntVar(&config.Metrics.Port, "metrics-port", 5003, "The port to use for serving metrics")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
