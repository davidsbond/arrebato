// Package acl provides all functionality within arrebato regarding access-control lists. This includes both gRPC, raft
// and data store interactions.
package acl

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"google.golang.org/protobuf/proto"

	"github.com/davidsbond/arrebato/internal/proto/arrebato/acl/v1"
)

// ErrInvalidACL is the error given when an ACL is invalid.
var ErrInvalidACL = errors.New("invalid ACL")

// Normalize the ACL, removing duplicate permissions and entries. This function will also detect invalid entries and
// return an error if one is found.
func Normalize(ctx context.Context, a *acl.ACL) (*acl.ACL, error) {
	normalizedACL := &acl.ACL{}

	for i, entry := range a.GetEntries() {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// We perform basic validation checks that contain the invalid entry's index, this way the user can easily
		// see which one is incorrect.
		switch {
		case entry.GetClient() == "":
			return nil, fmt.Errorf("%w: invalid entry %v: client name must be provided", ErrInvalidACL, i)
		case entry.GetTopic() == "":
			return nil, fmt.Errorf("%w: invalid entry %v: topic name must be provided", ErrInvalidACL, i)
		case len(entry.GetPermissions()) == 0:
			return nil, fmt.Errorf("%w: invalid entry %v: at least one permission is required", ErrInvalidACL, i)
		}

		flattenedPermissions := make(map[acl.Permission]struct{})

		// Having an unspecified permission is never allowed, so we check for that
		for _, permission := range entry.GetPermissions() {
			if permission != acl.Permission_PERMISSION_UNSPECIFIED {
				// Since permissions are just a slice, the same permission could be in the slice twice, we put these
				// into a map so that we don't have duplicates.
				flattenedPermissions[permission] = struct{}{}
				continue
			}

			return nil, fmt.Errorf("%w: invalid permission in entry %v: %s", ErrInvalidACL, i, permission)
		}

		permissions := make([]acl.Permission, 0)
		for permission := range flattenedPermissions {
			permissions = append(permissions, permission)
		}

		// Since we iterated over a map, ordering is non-deterministic, so we sort the permissions
		// here.
		sort.Slice(permissions, func(i, j int) bool {
			return permissions[i] < permissions[j]
		})

		newEntry := &acl.Entry{
			Topic:       entry.GetTopic(),
			Client:      entry.GetClient(),
			Permissions: permissions,
		}

		var duplicate bool
		// Because entries are a slice, the same entry could appear twice, so we check for those and don't append
		// if we've found a duplicate.
		for _, normalizedEntry := range normalizedACL.GetEntries() {
			if proto.Equal(entry, normalizedEntry) {
				duplicate = true
				break
			}
		}

		if duplicate {
			continue
		}

		normalizedACL.Entries = append(normalizedACL.Entries, newEntry)
	}

	return normalizedACL, nil
}
