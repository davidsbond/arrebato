package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	)

	return cmd
}

func listTopics() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "topics [flags]",
		Short: "List all topics",
		Long:  "This command returns a list of all topics within the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClient(cmd.Context())
			if err != nil {
				return err
			}

			topics, err := client.Topics(cmd.Context())
			if err != nil {
				return err
			}

			names := make([]string, len(topics))
			for i, topic := range topics {
				names[i] = topic.Name
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(names)
			}

			for _, name := range names {
				fmt.Println(name)
			}

			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&jsonOut, "json", false, "Output topic names in JSON format")

	return cmd
}
