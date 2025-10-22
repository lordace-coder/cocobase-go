package cocobase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) InitAuth(ctx context.Context) error {
	if c.storage == nil {
		return nil
	}

	token, err := c.storage.Get("cocobase-token")
	if err != nil {
		return nil
	}

	c.mu.Lock()
	c.token = token
	c.mu.Unlock()

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.user = user
	c.mu.Unlock()

	return nil
}

func (c *Client) Login(ctx context.Context, email, password string) error {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	
	resp, err := c.request(ctx, http.MethodPost, "/auth-collections/login", body, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := c.SetToken(tokenResp.AccessToken); err != nil {
		return err
	}

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.user = user
	c.mu.Unlock()

	return nil
}

func (c *Client) Register(ctx context.Context, email, password string, data map[string]interface{}) error {
	body := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	
	if data != nil {
		body["data"] = data
	}
	
	resp, err := c.request(ctx, http.MethodPost, "/auth-collections/signup", body, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := c.SetToken(tokenResp.AccessToken); err != nil {
		return err
	}

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.user = user
	c.mu.Unlock()

	return nil
}

func (c *Client) Logout() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.token = ""
	c.user = nil
	
	if c.storage != nil {
		return c.storage.Delete("cocobase-token")
	}
	
	return nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (*AppUser, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("user is not authenticated")
	}
	
	resp, err := c.request(ctx, http.MethodGet, "/auth-collections/user", nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user AppUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if c.storage != nil {
		userData, _ := json.Marshal(user)
		c.storage.Set("cocobase-user", string(userData))
	}

	return &user, nil
}

func (c *Client) UpdateUser(ctx context.Context, data map[string]interface{}, email, password *string) (*AppUser, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("user is not authenticated")
	}

	body := make(map[string]interface{})
	
	if data != nil {
		c.mu.RLock()
		currentData := make(map[string]interface{})
		if c.user != nil {
			currentData = c.user.Data
		}
		c.mu.RUnlock()
		
		merged := mergeData(currentData, data)
		body["data"] = merged
	}
	
	if email != nil {
		body["email"] = *email
	}
	
	if password != nil {
		body["password"] = *password
	}
	
	resp, err := c.request(ctx, http.MethodPatch, "/auth-collections/user", body, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user AppUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.mu.Lock()
	c.user = &user
	c.mu.Unlock()

	if c.storage != nil {
		userData, _ := json.Marshal(user)
		c.storage.Set("cocobase-user", string(userData))
	}

	return &user, nil
}

func mergeData(current, updates map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for k, v := range current {
		result[k] = v
	}
	
	for k, v := range updates {
		result[k] = v
	}
	
	return result
}
