# `baton-googleplay`

`baton-googleplay` is a connector for Google Play Console built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the [Google Play Developer API](https://developers.google.com/android-publisher) to sync data about users and their developer account permissions.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

## Prerequisites

To use this connector you need a **Google Cloud service account** with access to the Google Play Developer API:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create or select a project, then enable the **Google Play Android Developer API**
3. Go to **IAM & Admin** → **Service Accounts** → **Create Service Account**
4. Download the JSON key file for the service account
5. In [Google Play Console](https://play.google.com/console/) → **Users and permissions**, invite the service account email and grant appropriate permissions
6. Find your **Developer Account ID** in Play Console → **Settings** → **Developer account** → **Developer account ID**

## Getting Started

### source

```bash
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/EqualsGroup/baton-googleplay/cmd/baton-googleplay@main

baton-googleplay \
  --service-account-key-path "/path/to/service-account.json" \
  --developer-id "1234567890"

baton resources
```

### Environment variables

All flags can be set via environment variables:

```bash
export BATON_SERVICE_ACCOUNT_KEY_PATH="/path/to/service-account.json"
export BATON_DEVELOPER_ID="1234567890"

baton-googleplay
baton resources
```

## Data Model

`baton-googleplay` syncs the following resources:

| Resource Type | Description |
|--------------|-------------|
| Developer Account | The top-level Google Play developer account. Exposes one entitlement per developer-level permission. |
| User | Users with access to the developer account, including their email and permission grants. |

### Developer-Level Permissions

The following permissions are modeled as entitlements on the Developer Account:

| Permission | Description |
|-----------|-------------|
| `CAN_SEE_ALL_INFORMATION` | View app information and download bulk reports |
| `CAN_VIEW_FINANCIAL_DATA_GLOBAL` | View financial data, orders, and cancellation survey responses |
| `CAN_MANAGE_PERMISSIONS_GLOBAL` | Manage permissions for other users |
| `CAN_EDIT_GAMES_GLOBAL` | Edit Play Games Services settings |
| `CAN_PUBLISH_GAMES_GLOBAL` | Publish Play Games Services settings |
| `CAN_REPLY_TO_REVIEWS_GLOBAL` | Reply to user reviews |
| `CAN_MANAGE_PUBLIC_APKS_GLOBAL` | Manage app releases and configs |
| `CAN_MANAGE_TRACK_APKS_GLOBAL` | Manage testing tracks and edit apps |
| `CAN_MANAGE_TRACK_USERS_GLOBAL` | Manage lists of testers |
| `CAN_MANAGE_PUBLIC_LISTING_GLOBAL` | Manage store listing, pricing, and distribution |
| `CAN_MANAGE_DRAFT_APPS_GLOBAL` | Create and edit draft apps |
| `CAN_MANAGE_ORDERS_GLOBAL` | Manage orders and subscriptions |

## Provisioning

| Action | Resource | Description |
|--------|----------|-------------|
| Grant | Developer Account Permission | Add a developer-level permission to a user. |
| Revoke | Developer Account Permission | Remove a developer-level permission from a user. |

## Contributing, Support and Issues

We welcome contributions and ideas. If you have questions, problems, or ideas: please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

## Command Line Usage

```
baton-googleplay

Usage:
  baton-googleplay [flags]
  baton-googleplay [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --service-account-key-path string   required: Path to the Google service account JSON key file ($BATON_SERVICE_ACCOUNT_KEY_PATH)
      --developer-id string               required: Google Play Developer Account ID ($BATON_DEVELOPER_ID)
      --client-id string                  The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string              The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                       The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                              help for baton-googleplay
      --log-format string                 The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                  The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                      This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
  -v, --version                           version for baton-googleplay
```
