package arrebato

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	aclsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/acl/service/v1"
	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
)

type (
	// The ACL type represents the server's access-control list. It describes which clients are able to produce/consume
	// on desired topics.
	ACL struct {
		Entries []ACLEntry
	}

	// The ACLEntry type represents the relationship between a single client and topic.
	ACLEntry struct {
		// The name of the topic.
		Topic string
		// The client identifier the entry refers to. In an insecure environment, this can be an arbitrary string that the
		// client will use to identify itself in the request metadata. When using mutual TLS, this will be a SPIFFE ID that
		// the client will include in its TLS certificate.
		Client string
		// The permissions the client has on the topic.
		Permissions []ACLPermission
	}

	// The ACLPermission type is an enumeration that denotes an action a client can make on a topic.
	ACLPermission uint
)

// ErrNoACL is the error given when querying the server ACL before one has been created.
var ErrNoACL = errors.New("no acl")

// Constants for ACL permissions.
const (
	ACLPermissionUnspecified ACLPermission = iota
	ACLPermissionProduce
	ACLPermissionConsume
)

func (p ACLPermission) String() string {
	return acl.Permission_name[int32(p)]
}

func (a ACL) toProto() *acl.ACL {
	p := &acl.ACL{}

	for _, entry := range a.Entries {
		e := &acl.Entry{
			Topic:  entry.Topic,
			Client: entry.Client,
		}

		for _, permission := range entry.Permissions {
			e.Permissions = append(e.Permissions, acl.Permission(permission))
		}

		p.Entries = append(p.Entries, e)
	}

	return p
}

// SetACL updates the server's ACL.
func (c *Client) SetACL(ctx context.Context, acl ACL) error {
	svc := aclsvc.NewACLServiceClient(c.cluster.leader())
	_, err := svc.Set(ctx, &aclsvc.SetRequest{
		Acl: acl.toProto(),
	})
	switch {
	case status.Code(err) == codes.FailedPrecondition:
		c.cluster.findLeader(ctx)
		return c.SetACL(ctx, acl)
	default:
		return err
	}
}

// ACL returns the server's ACL. Returns ErrNoACL if an ACL has not been created.
func (c *Client) ACL(ctx context.Context) (ACL, error) {
	svc := aclsvc.NewACLServiceClient(c.cluster.any())
	resp, err := svc.Get(ctx, &aclsvc.GetRequest{})
	switch {
	case status.Code(err) == codes.NotFound:
		return ACL{}, ErrNoACL
	case err != nil:
		return ACL{}, err
	default:
		return aclFromProto(resp.GetAcl()), nil
	}
}

func aclFromProto(in *acl.ACL) ACL {
	var out ACL

	for _, entry := range in.GetEntries() {
		e := ACLEntry{
			Topic:  entry.GetTopic(),
			Client: entry.GetClient(),
		}

		for _, permission := range entry.GetPermissions() {
			e.Permissions = append(e.Permissions, ACLPermission(permission))
		}

		out.Entries = append(out.Entries, e)
	}

	return out
}
