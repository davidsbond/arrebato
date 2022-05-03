package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/pkg/arrebato"
)

// Delete returns a command for deleting various server resources.
func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete server resources (topics, etc.)",
		Long:  "This command is the parent command for deleting server resources",
	}

	cmd.AddCommand(
		deleteTopic(),
	)

	return cmd
}

func deleteTopic() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:               "topic [flags] <name>",
		Short:             "Deletes a topic",
		Long:              "This command deletes a topic from the server and all its messages",
		Args:              cobra.ExactValidArgs(1),
		ValidArgsFunction: validTopicNames,
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			if force {
				return client.DeleteTopic(ctx, args[0])
			}

			ok, err := confirmf(os.Stdin, os.Stdout, "Are you sure you want to delete topic %s?", args[0])
			switch {
			case err != nil:
				return err
			case !ok:
				return nil
			default:
				return client.DeleteTopic(ctx, args[0])
			}
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&force, "force", false, "Do not prompt for confirmation")

	return cmd
}
