package e2e_test

import (
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/testutil"
	"github.com/davidsbond/arrebato/pkg/arrebato"
)

type (
	ACLSuite struct {
		suite.Suite

		client *arrebato.Client
	}
)

func (s *ACLSuite) TestACLRulesApply() {
	ctx := testutil.Context(s.T())

	topic := arrebato.Topic{
		Name:                   "test-topic",
		MessageRetentionPeriod: time.Hour,
	}

	require.NoError(s.T(), s.client.CreateTopic(ctx, topic))

	s.Run("It should set the server ACL", func() {
		acl := arrebato.ACL{
			Entries: []arrebato.ACLEntry{
				{
					Topic:  "test-topic",
					Client: "test-client",
					Permissions: []arrebato.ACLPermission{
						arrebato.ACLPermissionProduce,
					},
				},
			},
		}

		require.NoError(s.T(), s.client.SetACL(ctx, acl))
	})

	s.Run("The ACL should allow us to produce messages", func() {
		producer := s.client.NewProducer("test-topic")
		require.NoError(s.T(), producer.Produce(ctx, arrebato.Message{
			Payload: timestamppb.Now(),
		}))
	})

	s.Run("As a different client, the ACL should not let us produce messages", func() {
		client, err := arrebato.Dial(ctx, arrebato.DefaultConfig(":5002"))
		require.NoError(s.T(), err)

		producer := client.NewProducer("test-topic")
		require.Error(s.T(), producer.Produce(ctx, arrebato.Message{
			Payload: timestamppb.Now(),
		}))
	})
}
