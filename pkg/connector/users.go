package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-googleplay/pkg/client"
)

type userBuilder struct {
	client *client.Client
}

func (b *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func newUserResource(user client.User, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	// Extract the display name from the user. The API may return a name field
	// or we fall back to the email.
	displayName := user.Email
	if user.Name != "" {
		displayName = user.Name
	}

	profile := map[string]interface{}{
		"access_state":                  user.AccessState,
		"developer_account_permissions": user.DeveloperAccountPermissions,
	}

	if user.ExpirationTime != "" {
		profile["expiration_time"] = user.ExpirationTime
	}

	if len(user.Grants) > 0 {
		grantInfo := make([]map[string]interface{}, 0, len(user.Grants))
		for _, g := range user.Grants {
			grantInfo = append(grantInfo, map[string]interface{}{
				"package_name":          g.PackageName,
				"app_level_permissions": g.AppLevelPermissions,
			})
		}
		profile["app_grants"] = grantInfo
	}

	opts := []resourceSdk.UserTraitOption{
		resourceSdk.WithEmail(user.Email, true),
		resourceSdk.WithUserProfile(profile),
	}

	var resourceOpts []resourceSdk.ResourceOption
	if parentResourceID != nil {
		resourceOpts = append(resourceOpts, resourceSdk.WithParentResourceID(parentResourceID))
	}

	return resourceSdk.NewUserResource(
		displayName,
		userResourceType,
		user.Email,
		opts,
		resourceOpts...,
	)
}

// List returns all users in the Google Play developer account.
func (b *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	var pageToken string
	if pToken != nil {
		pageToken = pToken.Token
	}

	usersResp, err := b.client.ListUsers(ctx, pageToken, 100)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-googleplay: failed to list users: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(usersResp.Users))
	for _, user := range usersResp.Users {
		resource, err := newUserResource(user, parentResourceID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-googleplay: failed to create resource for user %s: %w", user.Email, err)
		}
		resources = append(resources, resource)
	}

	return resources, usersResp.NextPageToken, nil, nil
}

// Entitlements returns an empty slice for users since they don't have their own entitlements.
func (b *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants returns an empty slice for users since grants are modeled on the developer account.
func (b *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(c *client.Client) *userBuilder {
	return &userBuilder{
		client: c,
	}
}
