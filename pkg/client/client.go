package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	baseURL = "https://androidpublisher.googleapis.com/androidpublisher/v3"
	scope   = "https://www.googleapis.com/auth/androidpublisher"
)

// Client is the Google Play Developer API client.
type Client struct {
	httpClient  *http.Client
	developerID string
}

// New creates a new Google Play API client using a service account key file.
func New(ctx context.Context, serviceAccountKeyPath, developerID string) (*Client, error) {
	keyData, err := os.ReadFile(serviceAccountKeyPath)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to read service account key file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(keyData, scope)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to parse service account key: %w", err)
	}

	httpClient := config.Client(ctx)

	return &Client{
		httpClient:  httpClient,
		developerID: developerID,
	}, nil
}

// ListUsers lists all users for the developer account.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users/list
func (c *Client) ListUsers(ctx context.Context, pageToken string, pageSize int) (*ListUsersResponse, error) {
	u := fmt.Sprintf("%s/developers/%s/users", baseURL, c.developerID)

	params := url.Values{}
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}

	u = u + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to create list users request: %w", err)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to list users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp)
	}

	var result ListUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to decode list users response: %w", err)
	}

	return &result, nil
}

// GetUser retrieves a specific user by email.
func (c *Client) GetUser(ctx context.Context, email string) (*User, error) {
	u := fmt.Sprintf("%s/developers/%s/users/%s", baseURL, c.developerID, url.PathEscape(email))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to create get user request: %w", err)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to decode get user response: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user with the given email and permissions.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users/create
func (c *Client) CreateUser(ctx context.Context, email string, permissions []string) (*User, error) {
	u := fmt.Sprintf("%s/developers/%s/users", baseURL, c.developerID)

	user := User{
		Email:                       email,
		DeveloperAccountPermissions: permissions,
	}

	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to marshal create user body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to create create-user request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to create user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.handleError(resp)
	}

	var created User
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to decode create user response: %w", err)
	}

	return &created, nil
}

// DeleteUser removes a user from the developer account.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users/delete
func (c *Client) DeleteUser(ctx context.Context, email string) error {
	u := fmt.Sprintf("%s/developers/%s/users/%s", baseURL, c.developerID, url.PathEscape(email))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return fmt.Errorf("baton-googleplay: failed to create delete user request: %w", err)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return fmt.Errorf("baton-googleplay: failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return c.handleError(resp)
	}

	return nil
}

// PatchUser updates a user's developer account permissions.
// https://developers.google.com/android-publisher/api-ref/rest/v3/users/patch
func (c *Client) PatchUser(ctx context.Context, email string, permissions []string) (*User, error) {
	u := fmt.Sprintf("%s/developers/%s/users/%s", baseURL, c.developerID, url.PathEscape(email))

	params := url.Values{}
	params.Set("updateMask", "developerAccountPermissions")

	u = u + "?" + params.Encode()

	patchBody := User{
		DeveloperAccountPermissions: permissions,
	}

	body, err := json.Marshal(patchBody)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to marshal patch user body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to create patch user request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to patch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp)
	}

	var updated User
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, fmt.Errorf("baton-googleplay: failed to decode patch user response: %w", err)
	}

	return &updated, nil
}

// Validate checks that the credentials and developer ID are valid by attempting to list users.
func (c *Client) Validate(ctx context.Context) error {
	_, err := c.ListUsers(ctx, "", -1)
	return err
}

// do executes an HTTP request and returns the response.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	l := ctxzap.Extract(ctx)
	l.Debug("making API request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Check for OAuth2 token errors.
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			return nil, fmt.Errorf("baton-googleplay: OAuth2 token error: %s", rErr.Body)
		}
		return nil, err
	}

	l.Debug("API response",
		zap.Int("status", resp.StatusCode),
	)

	return resp, nil
}

// handleError reads the response body and returns a descriptive error.
func (c *Client) handleError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("baton-googleplay: API error (status %d), failed to read body: %w", resp.StatusCode, err)
	}

	return fmt.Errorf("baton-googleplay: API error (status %d): %s", resp.StatusCode, string(body))
}
