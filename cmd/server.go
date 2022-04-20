package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/davidsbond/arrebato/internal/server"
)

// Server returns a cobra.Command that can be used to start an arrebato server instance.
func Server(version string) *cobra.Command {
	config := server.Config{
		Version: version,
	}

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
	flags.DurationVar(&config.Raft.Timeout, "raft-timeout", time.Minute, "The timeout to use for raft transport")
	flags.IntVar(&config.Raft.MaxSnapshots, "raft-max-snapshots", 3, "The maximum number of raft snapshots to store")
	flags.BoolVar(&config.Raft.NonVoter, "raft-non-voter", false, "If true, this server will never obtain leadership and will be a read-only replica")

	// gRPC configuration
	flags.IntVar(&config.GRPC.Port, "grpc-port", 5000, "The port to use for gRPC transport")

	// Serf configuration
	flags.IntVar(&config.Serf.Port, "serf-port", 5001, "The port to use for serf transport")
	flags.StringVar(&config.Serf.EncryptionKeyFile, "serf-encryption-key", "", "The path to the encryption key file used for serf transport")

	// TLS configuration
	flags.StringVar(&config.TLS.CertFile, "tls-cert", "", "The location of the TLS cert file for the server to use")
	flags.StringVar(&config.TLS.KeyFile, "tls-key", "", "The location of the TLS key file for the server to use")
	flags.StringVar(&config.TLS.CAFile, "tls-ca", "", "The location of the CA certificate that signs client certificates")

	// Prometheus configuration
	flags.IntVar(&config.Metrics.Port, "metrics-port", 5002, "The port to use for serving metrics")

	return cmd
}
