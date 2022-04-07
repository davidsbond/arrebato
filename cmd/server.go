package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/davidsbond/arrebato/internal/server"
)

// Server returns a cobra.Command that can be used to start an arrebato server instance.
func Server() *cobra.Command {
	var config server.Config

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start an arrebato server",
		Long:  "This command starts an arrebato server, either standalone or as part of an existing cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			undo, err := maxprocs.Set()
			if err != nil {
				return fmt.Errorf("failed to automatically set GOMAXPROCS: %w", err)
			}
			defer undo()

			svr, err := server.New(config)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			return svr.Start(cmd.Context())
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
	flags.DurationVar(&config.Raft.Timeout, "raft-timeout", time.Second, "The timeout to use for raft transport")
	flags.IntVar(&config.Raft.MaxSnapshots, "raft-max-snapshots", 3, "The maximum number of raft snapshots to store")
	flags.IntVar(&config.Raft.MaxPool, "raft-max-pool", 3, "The maximum number of connections in the raft connection pool")
	flags.BoolVar(&config.Raft.NonVoter, "raft-non-voter", false, "If true, this server will never obtain leadership and will be a read-only replica")

	// Serf configuration
	flags.IntVar(&config.Serf.Port, "serf-port", 5001, "The port to use for serf transport")
	flags.StringVar(&config.Serf.EncryptionKeyFile, "serf-encryption-key", "", "The path to the encryption key file used for serf transport")

	// gRPC configuration
	flags.IntVar(&config.GRPC.Port, "grpc-port", 5002, "The port to use for gRPC transport")
	flags.StringVar(&config.GRPC.TLSCertFile, "grpc-tls-cert", "", "The location of the TLS cert file for the server to use")
	flags.StringVar(&config.GRPC.TLSKeyFile, "grpc-tls-key", "", "The location of the TLS key file for the server to use")
	flags.StringVar(&config.GRPC.TLSCAFile, "grpc-tls-ca", "", "The location of the CA certificate that signs client certificates")

	// Prometheus configuration
	flags.IntVar(&config.Metrics.Port, "metrics-port", 5003, "The port to use for serving metrics")

	return cmd
}
