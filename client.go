package asl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	ClientID     string
	ClientSecret string

	HTTPClient http.Client

	BaseURL string
	Token   string
}

// Authenticate will grab a fresh JWT, replacing
// the existing cached token
func (c *Client) Authenticate() error {
	return fmt.Errorf("not implemented")
}

func (c *Client) makeReq(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string][]string, 2)
	headers["Content-Type"] = []string{"application/json"}

	if t := c.Token; t != "" {
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
