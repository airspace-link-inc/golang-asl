package asl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ClientID     string
	ClientSecret string

	HTTPClient http.Client

	BaseURL string
	Token
}

type Token struct {
	AccessToken string    `json:"accessToken"`
	Expiration  time.Time `json:"expires"`
	Scopes      string    `json:"scope"`
}

func (t Token) Expired() bool {
	return time.Now().After(t.Expiration)
}

// Authenticate will grab a fresh JWT, replacing
// the existing cached token
func (c *Client) Authenticate(ctx context.Context, scopes ...string) error {
	if c.ClientID == "" || c.ClientSecret == "" {
		return fmt.Errorf("missing client ID or client secret")
	}

	clientCredentials := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
	}

	if len(scopes) > 0 {
		clientCredentials["scope"] = []string{strings.Join(scopes, " ")}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/v1/oauth/token",
		strings.NewReader(clientCredentials.Encode()),
	)

	if err != nil {
		return err
	}

	auth0Resp, err := apiReq[Token](&c.HTTPClient, req)
	if err != nil {
		return err
	}

	if auth0Resp.Status != 200 {
		return fmt.Errorf("error authenticating: %s", auth0Resp.Msg)
	}

	c.Token = auth0Resp.Data
	return err
}

func (c *Client) makeReq(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string][]string, 2)
	headers["Content-Type"] = []string{"application/json"}

	if t := c.Token.AccessToken; t != "" {
		headers["Authorization"] = []string{"Bearer " + t}
	}

	req.Header = headers
	return req, nil
}

// apiReq will perform an HTTP request and then unmarshal the
// response into the target struct pointer
func apiReq[X any](client *http.Client, req *http.Request) (*Resp[X], error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp Resp[X]
	if err := json.Unmarshal(buf, &apiResp); err != nil {
		return nil, err
	}

	if code := resp.StatusCode; code >= 400 {
		return nil, &Err{Status: code, Msg: apiResp.Msg}
	}

	return &apiResp, nil
}
