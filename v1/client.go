package aslv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"sync"
)

var (
	// Production URL
	defalutBaseURL = "https://airhub-api.airspacelink.com"
	sandboxURL     = "https://airhub-api-sandbox.airspacelink.com"
)

type service struct {
	client *Client
}

type Client struct {
	// API endpoint
	// You can use the AirHub API in sandbox mode, which does not affect your live data or interact with production systems.
	// The base URL you use to make requests determines whether the request is live mode or sandbox mode.
	baseUrl *nurl.URL
	// API key issued by ASL to interact with API
	// Request access to the API by sending an email to developers@airspacelink.com.
	apiKey string
	client *http.Client

	// Auth token configs
	clientID     string
	clientSecret string

	// Cached token proptected by a mutex
	mu          sync.Mutex
	accessToken *Token

	// Services that Airspace Link provides
	Advisory *AdvisoryService
	Token    *TokenService
	Surface  *SurfaceService
	Check    *CheckService
}

func NewClient(httpClient *http.Client, baseURL, apiKey, clientID, clientSecret string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Attempt to parse input URL if it fails use default
	u, err := nurl.Parse(baseURL)
	if err != nil {
		u, _ = nurl.Parse(defalutBaseURL)
	}

	c := &Client{
		client:       httpClient,
		baseUrl:      u,
		apiKey:       apiKey,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	c.Advisory = &AdvisoryService{client: c}
	c.Token = &TokenService{client: c}
	c.Surface = &SurfaceService{client: c}
	c.Check = &CheckService{client: c}

	return c
}

func (c *Client) do(request *http.Request) (status int, body []byte, err error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return status, body, fmt.Errorf("airspacelink/v1: failed HTTP connection: %w", err)
	}

	status = resp.StatusCode

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("airspacelink/v1: failed to read resposne body: %w", err)
	}

	return status, body, err
}

// Checks the cached access token, if it satisfies the requirements to make the request it breaks.
// Otherwise it fetches a new token and replaces it in the cache.
func (c *Client) refreshToken(ctx context.Context, scope string) error {
	if c.accessToken != nil && c.accessToken.Valid() {
		return nil
	}

	// TODO: make sure token has the required scope

	t, err := c.Token.OauthToken(ctx, scope)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.accessToken = t

	return nil
}

// Reduce boilerplate code for creating a request with a JSON body
func (c *Client) requestWithPayload(ctx context.Context, method, url string, payload any, scope string) (*http.Request, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: failed to marshal request body to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: failed to create HTTP request: %w", err)
	}

	if err := c.refreshToken(ctx, scope); err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to retrieve access token: %w", err)
	}

	req.Header.Add("x-api-key", c.apiKey)
	req.Header.Add("Authorization", "Bearer "+c.accessToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *Client) post(ctx context.Context, url string, payload any, scope string) (status int, body []byte, err error) {
	req, err := c.requestWithPayload(ctx, http.MethodPost, url, payload, scope)
	if err != nil {
		return status, body, err
	}

	return c.do(req)
}

func (c *Client) delete(ctx context.Context, url string, scope string) (status int, body []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return status, body, fmt.Errorf("airspacelink/v1: failed to create HTTP request: %w", err)
	}

	if err := c.refreshToken(ctx, scope); err != nil {
		return status, body, fmt.Errorf("airspacelink/v1: unable to retrieve access token: %w", err)
	}

	req.Header.Add("x-api-key", c.apiKey)
	req.Header.Add("Authorization", "Bearer "+c.accessToken.AccessToken)

	return c.do(req)
}

func (c *Client) patch(ctx context.Context, url string, payload any, scope string) (status int, body []byte, err error) {
	req, err := c.requestWithPayload(ctx, http.MethodPatch, url, payload, scope)
	if err != nil {
		return status, body, err
	}

	return c.do(req)
}

// Checks the status code is 200 and then unmarshals the response body into the provided input
// Reduce boilerplate code for creating a request with a JSON body
func unmarshalStatusOk[T any](status int, body []byte, in T) error {
	if status != http.StatusOK {
		return fmt.Errorf("airspacelink/v1: status %d response: %s", status, body)
	}

	err := json.Unmarshal(body, &in)
	if err != nil {
		err = fmt.Errorf("airspacelink/v1: unable to unmarshal response: %w", err)
	}

	return err
}
