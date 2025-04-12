package authlete

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/pkg/config"
	"github.com/yamakenji24/golang-auth/pkg/logger"
)

type client struct {
	config     *config.Config
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

func NewClient(cfg *config.Config) *client {
	return &client{
		config:     cfg,
		httpClient: &http.Client{},
	}
}

func (c *client) sendRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &AuthleteError{Code: "REQUEST_ERROR", Message: "Failed to create request", Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AuthleteAccessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &AuthleteError{Code: "REQUEST_ERROR", Message: "Failed to send request", Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &AuthleteError{Code: "READ_ERROR", Message: "Failed to read response body", Err: err}
	}

	return respBody, nil
}

func (c *client) RequestAuthorization(params map[string]string) (*entity.AuthResponse, error) {
	apiURL := fmt.Sprintf("%s/%s/auth/authorization", c.config.AuthleteBaseURL, c.config.AuthleteServiceID)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	values.Set("redirect_uri", c.config.AuthleteRedirectURI)

	reqBody := map[string]string{
		"parameters": values.Encode(),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.LogError("Error marshaling request body: %v", err)
		return nil, &AuthleteError{Code: "MARSHAL_ERROR", Message: "Failed to marshal request body", Err: err}
	}

	body, err := c.sendRequest("POST", apiURL, jsonBody)
	if err != nil {
		logger.LogError("Error sending request: %v", err)
		return nil, err
	}

	var result entity.AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError("Error unmarshaling response body: %v", err)
		return nil, &AuthleteError{Code: "UNMARSHAL_ERROR", Message: "Failed to unmarshal response body", Err: err}
	}

	return &result, nil
}

func (c *client) IssueAuthorization(ticket string) (*entity.AuthResponse, error) {
	apiURL := fmt.Sprintf("%s/%s/auth/authorization/issue", c.config.AuthleteBaseURL, c.config.AuthleteServiceID)

	reqBody := map[string]string{
		"ticket":  ticket,
		"subject": "yamakenji",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.LogError("Error marshaling request body: %v", err)
		return nil, &AuthleteError{Code: "MARSHAL_ERROR", Message: "Failed to marshal request body", Err: err}
	}

	body, err := c.sendRequest("POST", apiURL, jsonBody)
	if err != nil {
		logger.LogError("Error sending request: %v", err)
		return nil, err
	}

	var result entity.AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError("Error unmarshaling response body: %v", err)
		return nil, &AuthleteError{Code: "UNMARSHAL_ERROR", Message: "Failed to unmarshal response body", Err: err}
	}

	return &result, nil
}

func (c *client) ExchangeToken(params map[string]string) (*entity.TokenResponse, error) {
	apiURL := fmt.Sprintf("%s/%s/auth/token", c.config.AuthleteBaseURL, c.config.AuthleteServiceID)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	reqBody := map[string]string{
		"parameters":   values.Encode(),
		"clientId":     c.config.AuthleteClientID,
		"clientSecret": c.config.AuthleteClientSecret,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.LogError("Error marshaling request body: %v", err)
		return nil, &AuthleteError{Code: "MARSHAL_ERROR", Message: "Failed to marshal request body", Err: err}
	}

	body, err := c.sendRequest("POST", apiURL, jsonBody)
	if err != nil {
		logger.LogError("Error sending request: %v", err)
		return nil, err
	}

	var result entity.TokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError("Error unmarshaling response body: %v", err)
		return nil, &AuthleteError{Code: "UNMARSHAL_ERROR", Message: "Failed to unmarshal response body", Err: err}
	}

	return &result, nil
}

// GetUserInfo アクセストークンからユーザー情報を取得
func (c *client) GetUserInfo(accessToken string) (entity.UserInfo, error) {
	apiURL := fmt.Sprintf("%s/%s/auth/userinfo", c.config.AuthleteBaseURL, c.config.AuthleteServiceID)
	fmt.Println("access_token: ", accessToken)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
			return entity.UserInfo{}, &AuthleteError{Code: "REQUEST_ERROR", Message: "Failed to create request", Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
			return entity.UserInfo{}, &AuthleteError{Code: "REQUEST_ERROR", Message: "Failed to send request", Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
			return entity.UserInfo{}, &AuthleteError{Code: "READ_ERROR", Message: "Failed to read response body", Err: err}
	}

	var userInfo entity.UserInfo
	if err := json.Unmarshal(respBody, &userInfo); err != nil {
			return entity.UserInfo{}, &AuthleteError{Code: "UNMARSHAL_ERROR", Message: "Failed to unmarshal response body", Err: err}
	}

	return userInfo, nil
}