package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ServiceAccountKeyPath = field.StringField(
		"service-account-key-path",
		field.WithDisplayName("Service Account Key Path"),
		field.WithDescription("Path to the Google service account JSON key file"),
		field.WithRequired(true),
	)

	DeveloperID = field.StringField(
		"developer-id",
		field.WithDisplayName("Developer ID"),
		field.WithDescription("Google Play Developer Account ID (numeric string)"),
		field.WithRequired(true),
	)

	ConfigurationFields = []field.SchemaField{
		ServiceAccountKeyPath,
		DeveloperID,
	}

	FieldRelationships = []field.SchemaFieldRelationship{}
)

var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Google Play Console"),
)
