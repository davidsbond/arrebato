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
		Raft: server.RaftConfig{
			Timeout:      time.Minute,
			MaxSnapshots: 3,
		},
		Metrics: server.MetricConfig{
			Port: 5002,
		},
		GRPC: server.GRPCConfig{
			Port: 5000,
		},
		Serf: server.SerfConfig{
			Port: 5001,
		},
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
	flags.BoolVar(&config.Raft.NonVoter, "read-only", false, "If true, this server will never obtain leadership and will be a read-only replica")

	// Serf configuration
	flags.StringVar(&config.Serf.EncryptionKeyFile, "serf-encryption-key", "", "The path to the encryption key file used for serf transport")

	// TLS configuration
	flags.StringVar(&config.TLS.CertFile, "tls-cert", "", "The location of the TLS cert file for the server to use")
	flags.StringVar(&config.TLS.KeyFile, "tls-key", "", "The location of the TLS key file for the server to use")
	flags.StringVar(&config.TLS.CAFile, "tls-ca", "", "The location of the CA certificate that signs client certificates")

	// Tracing configuration
	flags.BoolVar(&config.Tracing.Enabled, "tracing-enabled", false, "If true, enables OpenTelemetry span collection")
	flags.StringVar(&config.Tracing.Endpoint, "tracing-endpoint", "", "The HTTP(s) endpoint for the OpenTelemetry collector")

	return cmd
}
