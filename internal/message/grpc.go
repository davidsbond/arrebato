// Package message provides all aspects of the arrebato server for managing messages. Primarily the production and
// consumption of messages.
package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	"github.com/davidsbond/arrebato/internal/command"
	"github.com/davidsbond/arrebato/internal/consumer"
	"github.com/davidsbond/arrebato/internal/node"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
	consumerpb "github.com/davidsbond/arrebato/internal/proto/arrebato/consumer/v1"
	messagecmd "github.com/davidsbond/arrebato/internal/proto/arrebato/message/command/v1"
	messagesvc "github.com/davidsbond/arrebato/internal/proto/arrebato/message/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
	nodepb "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
	topicpb "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
	"github.com/davidsbond/arrebato/internal/signing"
	"github.com/davidsbond/arrebato/internal/topic"
)

type (
	// The GRPC type is a messagesvc.MessageServiceServer implementation that handles inbound gRPC requests to manage
	// and query Messages.
	GRPC struct {
		nodeName    string
		executor    Executor
		reader      Reader
		consumers   TopicIndexGetter
		publicKeys  PublicKeyGetter
		topics      TopicGetter
		topicOwners TopicOwnerGetter
		acl         ACL
	}

	// The Executor interface describes types that execute commands related to Message data.
	Executor interface {
		// Execute should perform actions corresponding to the provided command. The returned error should correspond
		// to the issue relevant to the command. For example, a command that creates a new Message should return
		// ErrNoTopic if the topic does not exist.
		Execute(ctx context.Context, cmd command.Command) error
	}

	// The Reader interface describes types that can read messages from a topic starting from a given index.
	Reader interface {
		// Read should start reading messages within a topic starting from a given index, invoking the ReadFunc for
		// each message.
		Read(ctx context.Context, topic string, startIndex uint64, fn ReadFunc) error
	}

	// The TopicIndexGetter interface describes types that can retrieve the current index on a topic for a consumer.
	TopicIndexGetter interface {
		// GetTopicIndex should return the number that describes the current index in the topic that a consumer has
		// read to. It should return ErrNoTopic if the topic does not exist.
		GetTopicIndex(ctx context.Context, topic string, consumerID string) (*consumerpb.TopicIndex, error)
	}

	// The ACL interface describes types that act as an access-control list to determine what clients are permitted
	// to do on certain topics.
	ACL interface {
		// Allowed should return true if the client has the given permission on a topic. If no ACL has been set up
		// within the server then it should always return true.
		Allowed(ctx context.Context, topic, client string, permission acl.Permission) (bool, error)
	}

	// The PublicKeyGetter interface describes types that can obtain a client's public signing key.
	PublicKeyGetter interface {
		// Get should return the client's public signing key. If there is no key then it should return signing.ErrNoPublicKey.
		Get(ctx context.Context, clientID string) ([]byte, error)
	}

	// The TopicGetter interface describes types that can obtain details on a specific topic.
	TopicGetter interface {
		// Get should return the named topic. It should return topic.ErrNoTopic if the topic does not exist.
		Get(ctx context.Context, name string) (*topicpb.Topic, error)
	}

	// The TopicOwnerGetter interface describes types that can query which node owns a topic.
	TopicOwnerGetter interface {
		// GetTopicOwner should return a node that is assigned to a given topic. Or return node.ErrNoNode if an owner
		// cannot be found for the given topic.
		GetTopicOwner(ctx context.Context, topicName string) (*nodepb.Node, error)
	}
)

// NewGRPC returns a new instance of the GRPC type that will modify Message data via commands sent to the Executor and
// read messages via the Reader implementation. The index of consumers will be obtained via the TopicIndexGetter implementation,
// permissions will be checked using the ACL implementation, topic details will be obtained via the TopicGetter implementation
// and client's public signing keys are obtained via the PublicKeyGetter implementation.
func NewGRPC(
	nodeName string,
	executor Executor,
	reader Reader,
	consumers TopicIndexGetter,
	acl ACL,
	publicKeys PublicKeyGetter,
	topics TopicGetter,
	owners TopicOwnerGetter,
) *GRPC {
	return &GRPC{
		nodeName:    nodeName,
		executor:    executor,
		reader:      reader,
		consumers:   consumers,
		acl:         acl,
		publicKeys:  publicKeys,
		topics:      topics,
		topicOwners: owners,
	}
}

// Register the GRPC service onto the grpc.ServiceRegistrar.
func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, health *health.Server) {
	messagesvc.RegisterMessageServiceServer(registrar, svr)
	health.SetServingStatus(messagesvc.MessageService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

// Produce a new message for a topic. Returns a codes.NotFound code if the topic does not exist.
func (svr *GRPC) Produce(ctx context.Context, request *messagesvc.ProduceRequest) (*messagesvc.ProduceResponse, error) {
	info := clientinfo.FromContext(ctx)
	allowed, err := svr.canProduce(ctx, info, request.GetMessage().GetTopic())
	switch {
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to check ACL: %v", err)
	case !allowed:
		return nil, status.Errorf(codes.PermissionDenied, "cannot produce messages on topic %s", request.GetMessage().GetTopic())
	}

	var verified bool
	if request.GetMessage().GetKey() != nil {
		verified, err = svr.verifySignature(ctx, info, request.GetMessage())
		switch {
		case errors.Is(err, signing.ErrNoPublicKey):
			// If there is no public key, we'll just mark the message as unverified
			break
		case err != nil:
			return nil, status.Errorf(codes.Internal, "failed to verify signature: %v", err)
		}
	}

	tp, err := svr.topics.Get(ctx, request.GetMessage().GetTopic())
	switch {
	case errors.Is(err, topic.ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetMessage().GetTopic())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to query topic: %v", err)
	case tp.GetRequireVerifiedMessages() && !verified:
		return nil, status.Errorf(codes.PermissionDenied, "topic %s only allows verified messages", request.GetMessage().GetTopic())
	}

	cmd := command.New(&messagecmd.CreateMessage{
		Message: &message.Message{
			Topic:     request.GetMessage().GetTopic(),
			Key:       request.GetMessage().GetKey(),
			Value:     request.GetMessage().GetValue(),
			Timestamp: timestamppb.Now(),
			Sender: &message.Sender{
				Id:           info.ID,
				KeySignature: request.GetMessage().GetSender().GetKeySignature(),
				Verified:     verified,
			},
		},
	})

	err = svr.executor.Execute(ctx, cmd)
	switch {
	case errors.Is(err, command.ErrNotLeader):
		return nil, status.Error(codes.FailedPrecondition, "this node is not the leader")
	case errors.Is(err, ErrNoTopic):
		return nil, status.Errorf(codes.NotFound, "topic %s does not exist", request.GetMessage().GetTopic())
	case err != nil:
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	default:
		return &messagesvc.ProduceResponse{}, nil
	}
}

// Consume messages from a topic. Returns a codes.NotFound code if the topic does not exist. In order to be allowed
// to consume the desired topic, the client must be given permission by the ACL, and must be calling the server that
// is assigned the desired topic.
func (svr *GRPC) Consume(request *messagesvc.ConsumeRequest, server messagesvc.MessageService_ConsumeServer) error {
	ctx := server.Context()
	info := clientinfo.FromContext(ctx)

	assigned, err := svr.isAssignedToTopic(ctx, request.GetTopic())
	switch {
	case errors.Is(err, node.ErrNoNode):
		return status.Errorf(codes.NotFound, "no owner was found for topic %s", request.GetTopic())
	case err != nil:
		return status.Errorf(codes.Internal, "failed to check topic assignments: %v", err)
	case !assigned:
		return status.Error(codes.PermissionDenied, "topics must be consumed via their assigned node")
	}

	allowed, err := svr.canConsume(ctx, info, request.GetTopic())
	switch {
	case err != nil:
		return status.Errorf(codes.Internal, "failed to check ACL: %v", err)
	case !allowed:
		return status.Errorf(codes.PermissionDenied, "cannot consume messages on topic %s", request.GetTopic())
	}

	topicIndex, err := svr.consumers.GetTopicIndex(ctx,
		request.GetTopic(),
		request.GetConsumerId(),
	)
	switch {
	case errors.Is(err, consumer.ErrNoTopic):
		return status.Errorf(codes.NotFound, "topic %s does not exist", request.GetTopic())
	case ctx.Err() != nil:
		return status.FromContextError(ctx.Err()).Err()
	case err != nil:
		return status.Errorf(codes.Internal, "failed to get last index for consumer: %v", err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		case <-ticker.C:
			// Messages here are read and forwarded to the client from the last index onwards. Once we reach the end
			// of all available messages we wait via the ticker and check again for new messages from the last known
			// index. This allows the client to stay connected and keep getting more messages, without over-polling
			// the data store.
			err = svr.reader.Read(ctx, request.GetTopic(), topicIndex.GetIndex(), func(ctx context.Context, m *message.Message) error {
				resp := &messagesvc.ConsumeResponse{Message: m}

				if err = server.Send(resp); err != nil {
					return fmt.Errorf("failed to send message: %w", err)
				}

				topicIndex.Index = m.GetIndex() + 1
				return nil
			})

			switch {
			case errors.Is(err, ErrNoTopic):
				return status.Errorf(codes.NotFound, "topic %s does not exist", request.GetTopic())
			case errors.Is(err, ErrNoMessages):
				// If no messages have been produced we'll just keep waiting and trying to read.
				continue
			case ctx.Err() != nil:
				return status.FromContextError(ctx.Err()).Err()
			case err != nil:
				return status.Errorf(codes.Internal, "failed to read messages: %v", err)
			}
		}
	}
}

func (svr *GRPC) canConsume(ctx context.Context, info clientinfo.ClientInfo, topic string) (bool, error) {
	return svr.acl.Allowed(ctx,
		topic,
		info.ID,
		acl.Permission_PERMISSION_CONSUME,
	)
}

func (svr *GRPC) canProduce(ctx context.Context, info clientinfo.ClientInfo, topic string) (bool, error) {
	return svr.acl.Allowed(ctx,
		topic,
		info.ID,
		acl.Permission_PERMISSION_PRODUCE,
	)
}

func (svr *GRPC) verifySignature(ctx context.Context, clientInfo clientinfo.ClientInfo, msg *message.Message) (bool, error) {
	if msg.GetSender() == nil {
		return false, nil
	}

	publicKey, err := svr.publicKeys.Get(ctx, clientInfo.ID)
	if err != nil {
		return false, err
	}

	return signing.Verify(msg.GetSender().GetKeySignature(), publicKey), nil
}

func (svr *GRPC) isAssignedToTopic(ctx context.Context, topicName string) (bool, error) {
	owner, err := svr.topicOwners.GetTopicOwner(ctx, topicName)
	if err != nil {
		return false, fmt.Errorf("failed to find topic owner: %w", err)
	}

	return owner.GetName() == svr.nodeName, nil
}
