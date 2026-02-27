package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-googleplay/pkg/client"
)

// Connector implements the baton connector interface for Google Play Console.
type Connector struct {
	client      *client.Client
	developerID string
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should
// be synced from the Google Play Console.
func (c *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newDeveloperAccountBuilder(c.client, c.developerID),
		newUserBuilder(c.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's
// authenticated HTTP client. Not used for this connector.
func (c *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (c *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Google Play Console",
		Description: "Connector for Google Play Console user and permission management via the Google Play Developer API.",
	}, nil
}

// Validate checks that the connector is properly configured by exercising the
// API credentials.
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	err := c.client.Validate(ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// New returns a new instance of the Google Play connector.
func New(ctx context.Context, serviceAccountKeyPath, developerID string) (*Connector, error) {
	c, err := client.New(ctx, serviceAccountKeyPath, developerID)
	if err != nil {
		return nil, err
	}
	return &Connector{
		client:      c,
		developerID: developerID,
	}, nil
}
