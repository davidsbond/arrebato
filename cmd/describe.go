package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	)

	return cmd
}

func describeTopic() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "topic [flags] <name>",
		Short: "Describe a topics",
		Long:  "This command returns information about a single topic",
		Args:  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClient(cmd.Context())
			if err != nil {
				return err
			}

			topic, err := client.Topic(cmd.Context(), args[0])
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
		},
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output topic information in JSON format")

	return cmd
}
