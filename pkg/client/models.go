package client

import "time"

// User represents a Google Play Console user.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users
type User struct {
	// The resource name of the user, e.g. "developers/{developerId}/users/{email}"
	Name string `json:"name"`

	// The user's email address.
	Email string `json:"email"`

	// The time at which the user's access expires, if set.
	// Uses RFC 3339 format. Empty string means no expiration.
	ExpirationTime string `json:"expirationTime,omitempty"`

	// Developer-level permissions granted to this user.
	DeveloperAccountPermissions []string `json:"developerAccountPermissions,omitempty"`

	// App-level access grants.
	AccessState string `json:"accessState,omitempty"`

	// Grants for specific apps.
	Grants []Grant `json:"grants,omitempty"`
}

// Grant represents app-level permissions for a user.
type Grant struct {
	// The resource name, e.g. "developers/{developerId}/users/{email}/grants/{packageName}"
	Name string `json:"name,omitempty"`

	// The package name of the app.
	PackageName string `json:"packageName"`

	// App-level permissions.
	AppLevelPermissions []string `json:"appLevelPermissions,omitempty"`
}

// ListUsersResponse is the response from listing users.
type ListUsersResponse struct {
	Users         []User `json:"users"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// ParseExpirationTime parses the expiration time string into a time.Time.
// Returns zero time if the string is empty or unparseable.
func (u *User) ParseExpirationTime() time.Time {
	if u.ExpirationTime == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, u.ExpirationTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

// DeveloperAccountPermission constants.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users#DeveloperAccountPermission
const (
	PermissionUnspecified                = "DEVELOPER_LEVEL_PERMISSION_UNSPECIFIED"
	PermissionCanSeeAllInformation       = "CAN_SEE_ALL_INFORMATION"
	PermissionCanViewFinancialData       = "CAN_VIEW_FINANCIAL_DATA_GLOBAL"
	PermissionCanManagePermissions       = "CAN_MANAGE_PERMISSIONS_GLOBAL"
	PermissionCanEditGamesGlobal         = "CAN_EDIT_GAMES_GLOBAL"
	PermissionCanPublishGamesGlobal      = "CAN_PUBLISH_GAMES_GLOBAL"
	PermissionCanReplyToReviewsGlobal    = "CAN_REPLY_TO_REVIEWS_GLOBAL"
	PermissionCanManagePublicListingGlob = "CAN_MANAGE_PUBLIC_LISTING_GLOBAL"
	PermissionCanManageTrackAPKsGlobal   = "CAN_MANAGE_TRACK_APKS_GLOBAL"
	PermissionCanManageTrackUsersGlobal  = "CAN_MANAGE_TRACK_USERS_GLOBAL"
	PermissionCanManagePublicAPKsGlobal  = "CAN_MANAGE_PUBLIC_APKS_GLOBAL"
	PermissionCanCreateManagedPlayApps   = "CAN_CREATE_MANAGED_PLAY_APPS_GLOBAL"
	PermissionCanManageOrdersGlobal      = "CAN_MANAGE_ORDERS_GLOBAL"
)

// AllDeveloperPermissions is the list of all known developer-level permissions.
var AllDeveloperPermissions = []string{
	PermissionCanSeeAllInformation,
	PermissionCanViewFinancialData,
	PermissionCanManagePermissions,
	PermissionCanEditGamesGlobal,
	PermissionCanPublishGamesGlobal,
	PermissionCanReplyToReviewsGlobal,
	PermissionCanManagePublicListingGlob,
	PermissionCanManageTrackAPKsGlobal,
	PermissionCanManageTrackUsersGlobal,
	PermissionCanManagePublicAPKsGlobal,
	PermissionCanCreateManagedPlayApps,
	PermissionCanManageOrdersGlobal,
}

// PermissionDescriptions provides human-readable descriptions for permissions.
var PermissionDescriptions = map[string]string{
	PermissionCanSeeAllInformation:       "Can see all information",
	PermissionCanViewFinancialData:       "Can view financial data globally",
	PermissionCanManagePermissions:       "Can manage permissions globally",
	PermissionCanEditGamesGlobal:         "Can edit games globally",
	PermissionCanPublishGamesGlobal:      "Can publish games globally",
	PermissionCanReplyToReviewsGlobal:    "Can reply to reviews globally",
	PermissionCanManagePublicListingGlob: "Can manage public listing globally",
	PermissionCanManageTrackAPKsGlobal:   "Can manage track APKs globally",
	PermissionCanManageTrackUsersGlobal:  "Can manage track users globally",
	PermissionCanManagePublicAPKsGlobal:  "Can manage public APKs globally",
	PermissionCanCreateManagedPlayApps:   "Can create managed Play apps globally",
	PermissionCanManageOrdersGlobal:      "Can manage orders globally",
}
