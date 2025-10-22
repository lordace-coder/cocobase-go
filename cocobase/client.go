package cocobase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: DefaultTimeout,
		}
	}

	return &Client{
		baseURL:    strings.TrimSuffix(config.BaseURL, "/"),
		apiKey:     config.APIKey,
		httpClient: config.HTTPClient,
		storage:    config.Storage,
	}
}

func (c *Client) SetToken(token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.token = token
	
	if c.storage != nil {
		return c.storage.Set("cocobase-token", token)
	}
	
	return nil
}

func (c *Client) GetToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token != ""
}

func (c *Client) HasRole(role string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if c.user == nil {
		return false
	}
	
	for _, r := range c.user.Roles {
		if r == role {
			return true
		}
	}
	
	return false
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}, useDataKey bool) (*http.Response, error) {
	url := c.baseURL + path
	
	var bodyReader io.Reader
	if body != nil {
		var data interface{}
		if useDataKey {
			data = map[string]interface{}{"data": body}
		} else {
			data = body
		}
		
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", ContentTypeJSON)
	
	if c.apiKey != "" {
		req.Header.Set(HeaderAPIKey, c.apiKey)
	}
	
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()
	
	if token != "" {
		req.Header.Set(HeaderAuthorization, "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Method:     method,
			URL:        url,
			Body:       string(bodyBytes),
			Suggestion: getErrorSuggestion(resp.StatusCode, method),
		}
	}

	return resp, nil
}
