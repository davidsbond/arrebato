package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type (
	SigningSuite struct {
		suite.Suite

		client *arrebato.Client
	}
)

func (s *SigningSuite) TestKeys() {
	ctx := testutil.Context(s.T())

	var keyPair arrebato.KeyPair
	var err error

	s.T().Run("It should create a key pair for signing messages", func(t *testing.T) {
		keyPair, err = s.client.CreateSigningKeyPair(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, keyPair.PublicKey)
		require.NotEmpty(t, keyPair.PrivateKey)
	})

	s.T().Run("It should obtain the same public key from the server", func(t *testing.T) {
		publicKey, err := s.client.SigningPublicKey(ctx, "test-client")
		require.NoError(t, err)
		require.EqualValues(t, keyPair.PublicKey, publicKey)
	})

	// Create a topic to publish a signed message on
	require.NoError(s.T(), s.client.CreateTopic(ctx, arrebato.Topic{
		Name:                    "signing-topic",
		MessageRetentionPeriod:  time.Hour,
		ConsumerRetentionPeriod: time.Hour,
		RequireVerifiedMessages: true,
	}))

	// We need a fresh client so we can set the signing key
	client, err := arrebato.Dial(ctx, arrebato.Config{
		Addresses:         []string{":5002"},
		ClientID:          "test-client",
		MessageSigningKey: keyPair.PrivateKey,
	})
	require.NoError(s.T(), err)

	s.T().Run("It should sign the message if we provide a key", func(t *testing.T) {
		producer := client.NewProducer("signing-topic")
		require.NoError(t, producer.Produce(ctx, arrebato.Message{
			Value: structpb.NewStringValue("world"),
			Key:   structpb.NewStringValue("hello"),
		}))

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		consumer, err := client.NewConsumer(ctx, arrebato.ConsumerConfig{
			Topic:      "signing-topic",
			ConsumerID: "test-consumer",
		})

		require.NoError(t, err)
		err = consumer.Consume(ctx, func(ctx context.Context, m arrebato.Message) error {
			defer cancel()

			assert.True(t, m.Sender.Verified)
			assert.NotEmpty(t, m.Sender.KeySignature)
			assert.True(t, signing.Verify(m.Sender.KeySignature, keyPair.PublicKey))
			return nil
		})
	})
}
