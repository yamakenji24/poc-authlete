package authlete

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yamakenji24/golang-auth/pkg/config"
)

type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func TestRequestAuthorization(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ticket": "test-ticket"}`)),
		},
	}

	cfg := &config.Config{
		AuthleteBaseURL:     "http://test-server",
		AuthleteServiceID:   "test-service",
		AuthleteAccessToken: "test-token",
		AuthleteRedirectURI: "http://localhost:8081/auth/callback",
	}
	client := NewClient(cfg)
	client.httpClient = mockClient

	params := map[string]string{
		"response_type": "code",
		"client_id":     "test-client",
		"scope":         "openid",
	}

	resp, err := client.RequestAuthorization(params)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-ticket", resp.Ticket)
}

func TestIssueAuthorization(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ticket": "test-ticket"}`)),
		},
	}

	cfg := &config.Config{
		AuthleteBaseURL:     "http://test-server",
		AuthleteServiceID:   "test-service",
		AuthleteAccessToken: "test-token",
	}
	client := NewClient(cfg)
	client.httpClient = mockClient

	resp, err := client.IssueAuthorization("test-ticket")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-ticket", resp.Ticket)
}

func TestExchangeToken(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"accessToken": "test-access-token"}`)),
		},
	}

	cfg := &config.Config{
		AuthleteBaseURL:      "http://test-server",
		AuthleteServiceID:    "test-service",
		AuthleteAccessToken:  "test-token",
		AuthleteClientID:     "test-client",
		AuthleteClientSecret: "test-secret",
	}
	client := NewClient(cfg)
	client.httpClient = mockClient

	params := map[string]string{
		"grant_type":   "authorization_code",
		"code":         "test-code",
		"redirect_uri": "http://localhost:8081/auth/callback",
	}

	resp, err := client.ExchangeToken(params)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-access-token", resp.AccessToken)
}
