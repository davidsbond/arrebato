package cmd

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/davidsbond/arrebato/pkg/arrebato"
)

// Create returns a cobra.Command that is a parent command for all create operations. For example, its child commands
// can create topics, signing keys etc.
func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a server resource (topics, signing keys etc.)",
		Long:  "This command is the parent command for creating a server resource",
	}

	cmd.AddCommand(
		createTopic(),
		createSigningKey(),
	)

	return cmd
}

func createTopic() *cobra.Command {
	var (
		messageRetentionPeriod  time.Duration
		consumerRetentionPeriod time.Duration
		requireVerifiedMessages bool
	)

	cmd := &cobra.Command{
		Use:   "topic [flags] <name>",
		Args:  cobra.ExactValidArgs(1),
		Short: "Create a new topic",
		Long:  "This command creates a new topic with the configuration provided by the CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClient(cmd.Context())
			if err != nil {
				return err
			}

			name := args[0]
			return client.CreateTopic(cmd.Context(), arrebato.Topic{
				Name:                    name,
				MessageRetentionPeriod:  messageRetentionPeriod,
				ConsumerRetentionPeriod: consumerRetentionPeriod,
				RequireVerifiedMessages: requireVerifiedMessages,
			})
		},
	}

	flags := cmd.PersistentFlags()
	flags.DurationVar(&messageRetentionPeriod, "message-retention", 0, "The amount of time to store a message on a topic, zero meaning infinite retention")
	flags.DurationVar(&consumerRetentionPeriod, "consumer-retention", 0, "The amount of time to store a consumer's index on a topic, zero meaning infinite retention")
	flags.BoolVar(&requireVerifiedMessages, "require-verified-messages", false, "If true, messages will be rejected if their signature is not verified")

	return cmd
}

func createSigningKey() *cobra.Command {
	return &cobra.Command{
		Use:   "signing-key",
		Short: "Create a new signing key pair",
		Long:  "This command creates a new signing key pair for this client to use when producing messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClient(cmd.Context())
			if err != nil {
				return err
			}

			keyPair, err := client.CreateSigningKeyPair(cmd.Context())
			if err != nil {
				return err
			}

			fmt.Printf("Keys below are base64 encoded, they should be provided to the server decoded.\n\n")
			fmt.Printf("Public key:\t%s\n", base64.StdEncoding.EncodeToString(keyPair.PublicKey))
			fmt.Printf("Private key:\t%s\n", base64.StdEncoding.EncodeToString(keyPair.PrivateKey))
			return nil
		},
	}
}
