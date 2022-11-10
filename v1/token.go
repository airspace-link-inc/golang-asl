package aslv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenService service

type Token struct {
	AccessToken string    `json:"accessToken"`
	Expires     time.Time `json:"expires"`
	Scope       string    `json:"scope"`
}

// Returns true if the access token is non-empty and the token has not expired
func (t *Token) Valid() bool {
	return t.AccessToken != "" && time.Now().Before(t.Expires)
}

func (ts *TokenService) OauthToken(ctx context.Context, scopes ...string) (*Token, error) {
	params := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ts.client.clientID},
		"client_secret": {ts.client.clientSecret},
		"scope":         {strings.Join(scopes, " ")},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.client.baseUrl.String()+"/v1/oauth/token", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: failed to create HTTP request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("x-api-key", ts.client.apiKey)

	status, body, err := ts.client.do(req)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: failed HTTP connection: %w", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("airspacelink/v1: status %d response: %s", status, body)
	}

	var details struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		Data       struct {
			AccessToken string `json:"accessToken"`
			Expires     string `json:"expires"`
			Scope       string `json:"scope"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &details)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to unmarshal token response: %w", err)
	}

	expiration, err := time.Parse(time.RFC3339, details.Data.Expires)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to parse token expiration: %w", err)
	}

	return &Token{
		AccessToken: details.Data.AccessToken,
		Scope:       details.Data.Scope,
		Expires:     expiration,
	}, nil
}
