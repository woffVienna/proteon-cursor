package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/domain"
)

// Client calls the identity service to issue backoffice tokens.
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new identity client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type backofficeTokenRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	SubjectType string    `json:"subject_type"`
	TenantID    *string   `json:"tenant_id,omitempty"`
	Audience    string    `json:"audience"`
}

type backofficeTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// IssueBackofficeToken calls the identity internal endpoint and returns the result.
func (c *Client) IssueBackofficeToken(ctx context.Context, userID, subjectType, tenant string) (domain.LoginResult, error) {
	u, err := uuid.Parse(userID)
	if err != nil {
		return domain.LoginResult{}, fmt.Errorf("parse user id: %w", err)
	}

	reqBody := backofficeTokenRequest{
		UserID:      u,
		SubjectType: subjectType,
		Audience:    "backoffice",
	}
	if tenant != "" {
		reqBody.TenantID = &tenant
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&reqBody); err != nil {
		return domain.LoginResult{}, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/v1/backoffice-tokens", &buf)
	if err != nil {
		return domain.LoginResult{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return domain.LoginResult{}, fmt.Errorf("call identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.LoginResult{}, fmt.Errorf("identity returned status %d", resp.StatusCode)
	}

	var out backofficeTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.LoginResult{}, fmt.Errorf("decode response: %w", err)
	}

	return domain.LoginResult{
		AccessToken: out.AccessToken,
		ExpiresIn:   out.ExpiresIn,
	}, nil
}
