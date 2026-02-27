package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// developerAccountResourceType represents the Google Play Developer Account.
// This is the top-level resource that contains users and their permissions.
var developerAccountResourceType = &v2.ResourceType{
	Id:          "developer_account",
	DisplayName: "Developer Account",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

// userResourceType represents a user within the Google Play Console.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}
