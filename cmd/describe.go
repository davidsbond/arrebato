package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/pkg/arrebato"
)

// Describe returns a command for describing various server resources.
func Describe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe server resources (topics, etc.)",
		Long:  "This command is the parent command for describing server resources",
	}

	cmd.AddCommand(
		describeTopic(),
		describeNode(),
	)

	return cmd
}

func describeTopic() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "topic [flags] <name>",
		Short: "Describe a topic",
		Long:  "This command returns information about a single topic",
		Args:  cobra.ExactValidArgs(1),
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			topic, err := client.Topic(ctx, args[0])
			if err != nil {
				return err
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(topic)
			}

			fmt.Println("Name:", topic.Name)
			fmt.Println("Message Retention:", topic.MessageRetentionPeriod.String())
			fmt.Println("Consumer Retention:", topic.ConsumerRetentionPeriod.String())
			fmt.Println("Require Verified Messages:", topic.RequireVerifiedMessages)

			return nil
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output topic information in JSON format")

	return cmd
}

func describeNode() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "node [flags] <name>",
		Short: "Describe a node",
		Long:  "This command returns information about a single node in the cluster",
		Args:  cobra.ExactValidArgs(1),
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			node, err := client.Node(ctx, args[0])
			if err != nil {
				return err
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(node)
			}

			fmt.Println("Name:", node.Name)
			fmt.Println("Leader:", node.Leader)
			fmt.Println("Version:", node.Version)
			fmt.Println("Peers:", len(node.Peers))

			return nil
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output node information in JSON format")

	return cmd
}
