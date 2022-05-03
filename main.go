package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/cmd"
)

var version = "dev"

func main() {
	command := &cobra.Command{
		Use:     "arrebato",
		Version: version,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Long: "The command-line tool for managing arrebato servers.\n\n" +
			"Environment variables:\n" +
			"  ARREBATO_SERVERS         A comma-separated list of server addresses\n" +
			"  ARREBATO_CLIENT_ID       A client-identifier for calls made by this host\n" +
			"  ARREBATO_TLS             Indicates that TLS is used to communicate with the server\n" +
			"  ARREBATO_TLS_CA          The location of the CA certificate that signed the server's certificate\n" +
			"  ARREBATO_TLS_CERTIFICATE The location of the client certificate\n" +
			"  ARREBATO_TLS_KEY         The location of the client private key",
	}

	command.SetErr(os.Stderr)
	command.AddCommand(
		cmd.Server(version),
		cmd.Create(),
		cmd.List(),
		cmd.Describe(),
		cmd.Backup(),
		cmd.Delete(),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	if err := command.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}
