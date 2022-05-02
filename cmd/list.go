package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/internal/table"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

// List returns a command for listing various server resources.
func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List server resources (topics, etc.)",
		Long:  "This command is the parent command for listing server resources",
	}

	cmd.AddCommand(
		listTopics(),
		listSigningKeys(),
		listNodes(),
	)

	return cmd
}

func listTopics() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "topics [flags]",
		Short: "List all topics",
		Long:  "This command returns a list of all topics within the server",
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			topics, err := client.Topics(ctx)
			if err != nil {
				return err
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(topics)
			}

			builder := table.NewBuilder("NAME", "MESSAGE RETENTION PERIOD", "CONSUMER RETENTION PERIOD", "REQUIRE VERIFIED MESSAGES")
			for _, t := range topics {
				builder.AddRow(t.Name, t.MessageRetentionPeriod, t.ConsumerRetentionPeriod, t.RequireVerifiedMessages)
			}

			return builder.Build(ctx, os.Stdout)
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output topics in JSON format")

	return cmd
}

func listSigningKeys() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "signing-keys [flags]",
		Short: "List all public signing keys",
		Long:  "This command returns a list of all public signing keys within the server",
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			keys, err := client.PublicKeys(ctx)
			if err != nil {
				return err
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(keys)
			}

			builder := table.NewBuilder("CLIENT ID", "PUBLIC KEY")
			for _, k := range keys {
				builder.AddRow(k.ClientID, base64.StdEncoding.EncodeToString(k.PublicKey))
			}

			return builder.Build(ctx, os.Stdout)
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output public keys in JSON format")

	return cmd
}

func listNodes() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "nodes [flags]",
		Short: "List all nodes in the cluster",
		Long:  "This command returns a list of all nodes within the arrebato cluster",
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			nodes, err := client.Nodes(ctx)
			if err != nil {
				return err
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(nodes)
			}

			builder := table.NewBuilder("NAME", "LEADER", "VERSION", "TOPICS")
			for _, n := range nodes {
				builder.AddRow(n.Name, n.Leader, n.Version, len(n.Topics))
			}

			return builder.Build(ctx, os.Stdout)
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output nodes in JSON format")

	return cmd
}
