package connector

import (
	"context"
	"fmt"
	"slices"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-googleplay/pkg/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type developerAccountBuilder struct {
	client      *client.Client
	developerID string
}

func (b *developerAccountBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return developerAccountResourceType
}

// List returns the single developer account as a resource.
// There is exactly one developer account per connector instance.
func (b *developerAccountBuilder) List(_ context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	profile := map[string]interface{}{
		"developer_id": b.developerID,
	}

	resource, err := resourceSdk.NewGroupResource(
		fmt.Sprintf("Google Play Developer Account (%s)", b.developerID),
		developerAccountResourceType,
		b.developerID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
		resourceSdk.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: userResourceType.Id},
		),
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-googleplay: failed to create developer account resource: %w", err)
	}

	return []*v2.Resource{resource}, "", nil, nil
}

// Entitlements returns one permission entitlement per developer-level permission.
// These are the account-level permissions that can be granted to users.
func (b *developerAccountBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlements := make([]*v2.Entitlement, 0, len(client.AllDeveloperPermissions))

	for _, perm := range client.AllDeveloperPermissions {
		desc := client.PermissionDescriptions[perm]
		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(
			resource,
			perm,
			entitlement.WithDescription(fmt.Sprintf("%s: %s", perm, desc)),
			entitlement.WithDisplayName(fmt.Sprintf("%s %s", resource.DisplayName, perm)),
			entitlement.WithGrantableTo(userResourceType),
		))
	}

	return entitlements, "", nil, nil
}

// Grants returns a grant for each user's developer-level permissions.
// Each user may have multiple developer-level permissions, resulting in multiple grants.
func (b *developerAccountBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	var pageToken string
	if pToken != nil {
		pageToken = pToken.Token
	}

	usersResp, err := b.client.ListUsers(ctx, pageToken, -1)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-googleplay: failed to list users for grants: %w", err)
	}

	var grants []*v2.Grant

	for _, user := range usersResp.Users {
		userResourceID, err := resourceSdk.NewResourceID(userResourceType, user.Email)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-googleplay: failed to create resource ID for user %s: %w", user.Email, err)
		}

		for _, perm := range user.DeveloperAccountPermissions {
			if perm == client.PermissionUnspecified {
				continue
			}

			// Verify this is a known permission.
			if !slices.Contains(client.AllDeveloperPermissions, perm) {
				l.Debug("unknown developer permission, skipping grant",
					zap.String("permission", perm),
					zap.String("user", user.Email),
				)
				continue
			}

			grants = append(grants, grant.NewGrant(resource, perm, userResourceID))
		}
	}

	return grants, usersResp.NextPageToken, nil, nil
}

// Grant adds a developer-level permission to a user.
func (b *developerAccountBuilder) Grant(ctx context.Context, principal *v2.Resource, ent *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != userResourceType.Id {
		return nil, fmt.Errorf("baton-googleplay: expected principal to be a user, got %s", principal.Id.ResourceType)
	}

	email := principal.Id.Resource
	permission := ent.Slug

	// Get current user to check existing permissions.
	user, err := b.client.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to get user %s: %w", email, err)
	}

	// Check if user already has this permission.
	if slices.Contains(user.DeveloperAccountPermissions, permission) {
		return annotations.New(&v2.GrantAlreadyExists{}), nil
	}

	// Add the new permission to existing ones.
	newPerms := append(user.DeveloperAccountPermissions, permission)

	l.Debug("granting developer permission",
		zap.String("email", email),
		zap.String("permission", permission),
	)

	_, err = b.client.PatchUser(ctx, email, newPerms)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to grant permission %s to user %s: %w", permission, email, err)
	}

	return nil, nil
}

// Revoke removes a developer-level permission from a user.
func (b *developerAccountBuilder) Revoke(ctx context.Context, gnt *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	email := gnt.Principal.Id.Resource
	permission := gnt.Entitlement.Slug

	// Get current user to check existing permissions.
	user, err := b.client.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to get user %s: %w", email, err)
	}

	// Check if user has this permission.
	if !slices.Contains(user.DeveloperAccountPermissions, permission) {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	// Remove the permission.
	newPerms := make([]string, 0, len(user.DeveloperAccountPermissions))
	for _, p := range user.DeveloperAccountPermissions {
		if p != permission {
			newPerms = append(newPerms, p)
		}
	}

	l.Debug("revoking developer permission",
		zap.String("email", email),
		zap.String("permission", permission),
	)

	_, err = b.client.PatchUser(ctx, email, newPerms)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to revoke permission %s from user %s: %w", permission, email, err)
	}

	return nil, nil
}

func newDeveloperAccountBuilder(c *client.Client, developerID string) *developerAccountBuilder {
	return &developerAccountBuilder{
		client:      c,
		developerID: developerID,
	}
}
