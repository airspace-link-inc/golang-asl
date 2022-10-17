package asl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(mainTest *testing.T) {
	testCases := []struct {
		name        string
		client      Client
		ctx         context.Context
		scopes      []string
		expected    Token
		expectedErr string
	}{
		{
			name:        "missing client ID",
			expectedErr: "missing client ID or client secret",
		},
		{
			name: "missing client secret",
			client: Client{
				ClientID: "something",
			},
			expectedErr: "missing client ID or client secret",
		},
		{
			name: "canceled context",
			client: Client{
				ClientID:     "something",
				ClientSecret: "something",
				BaseURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"message": "error never should happen"}`))
				})).URL,
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectedErr: "context canceled",
		},
		{
			name: "HTTP 400",
			client: Client{
				ClientID:     "something",
				ClientSecret: "something",
				BaseURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(400)
					w.Write([]byte(`{"message": "something invalid happened"}`))
				})).URL,
			},
			expectedErr: "something invalid happened",
		},
		{
			name: "HTTP 500",
			client: Client{
				ClientID:     "something",
				ClientSecret: "something",
				BaseURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
					w.Write([]byte(`{"message": "something invalid happened (500)"}`))
				})).URL,
			},
			expectedErr: "something invalid happened (500)",
		},
		{
			name: "mock success",
			client: Client{
				ClientID:     "something",
				ClientSecret: "something",
				BaseURL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					w.Write([]byte(`{
						"statusCode": 200,
						"message": "success",
						"data": {
						  "accessToken": "xyz123",
						  "expires": "1000-02-03T07:43:22.987Z",
						  "scope": "airhub-api/advisory.read airhub-api/token.create"
						}
					  }`))
				})).URL,
			},
			expected: Token{
				AccessToken: "xyz123",
				Expiration:  time.Date(1000, 2, 3, 7, 43, 22, 987000000, time.UTC),
				Scopes:      "airhub-api/advisory.read airhub-api/token.create",
			},
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		if tc.ctx == nil {
			tc.ctx = context.Background()
		}

		actualErr := tc.client.Authenticate(tc.ctx)

		if tc.expectedErr != "" {
			t.Contains(actualErr.Error(), tc.expectedErr, tc.name)
			continue
		}

		if t.Nil(actualErr) {
			t.Equal(tc.expected, tc.client.Token, tc.name)
		}
	}
}
