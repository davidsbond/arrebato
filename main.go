package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/cmd"
)

var version string

func main() {
	command := &cobra.Command{
		Use:     "arrebato",
		Version: version,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
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

	command.AddCommand(
		cmd.Server(),
		cmd.Create(),
		cmd.List(),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	if err := command.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}
