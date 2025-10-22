#!/bin/bash

# Cocobase Go Client - Project Setup Script
# Run this from inside the cocobase-go folder
# Usage: chmod +x setup.sh && ./setup.sh

set -e

echo "🚀 Setting up Cocobase Go Client project structure..."

# Create directory structure
echo "📁 Creating directories..."
mkdir -p cocobase
mkdir -p examples/{basic,advanced,realtime,auth}
mkdir -p storage
mkdir -p tests

# Create go.mod
echo "📦 Creating go.mod..."
cat > go.mod << 'EOF'
module github.com/yourusername/cocobase-go

go 1.21

require (
	github.com/gorilla/websocket v1.5.1
)
EOF

# Create main package files
echo "📝 Creating cocobase package files..."

# types.go
cat > cocobase/types.go << 'EOF'
package cocobase

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	DefaultBaseURL          = "https://api.cocobase.com"
	DefaultTimeout          = 30 * time.Second
	ContentTypeJSON         = "application/json"
	HeaderAPIKey           = "x-api-key"
	HeaderAuthorization    = "Authorization"
)

type Client struct {
	baseURL    string
	apiKey     string
	token      string
	user       *AppUser
	httpClient *http.Client
	mu         sync.RWMutex
	storage    Storage
}

type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Storage    Storage
}

type Storage interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type Document struct {
	ID         string                 `json:"id"`
	Collection string                 `json:"collection"`
	Data       map[string]interface{} `json:"data"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type AppUser struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Roles     []string               `json:"roles"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type Connection struct {
	conn   *websocket.Conn
	name   string
	closed bool
	mu     sync.Mutex
}

type Event struct {
	Event string   `json:"event"`
	Data  Document `json:"data"`
}
EOF

# client.go
cat > cocobase/client.go << 'EOF'
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
EOF

# errors.go
cat > cocobase/errors.go << 'EOF'
package cocobase

import "fmt"

type APIError struct {
	StatusCode int
	Method     string
	URL        string
	Body       string
	Suggestion string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed: %s %s (status: %d)\nBody: %s\nSuggestion: %s",
		e.Method, e.URL, e.StatusCode, e.Body, e.Suggestion)
}

func getErrorSuggestion(status int, method string) string {
	switch status {
	case 401:
		return "Check if your API key is valid and properly set"
	case 403:
		return "You don't have permission to perform this action. Verify your access rights"
	case 404:
		return "The requested resource was not found. Verify the path and ID are correct"
	case 405:
		return fmt.Sprintf("The %s method is not allowed for this endpoint. Check the API documentation for supported methods", method)
	case 429:
		return "You've exceeded the rate limit. Please wait before making more requests"
	default:
		return "Check the API documentation and verify your request format"
	}
}
EOF

# query.go
cat > cocobase/query.go << 'EOF'
package cocobase

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryBuilder struct {
	filters    map[string]string
	orFilters  map[string][]string
	limit      int
	offset     int
	sort       string
	order      string
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{
		filters:   make(map[string]string),
		orFilters: make(map[string][]string),
	}
}

func (qb *QueryBuilder) Filter(field, operator string, value interface{}) *QueryBuilder {
	key := field
	if operator != "" && operator != "eq" {
		key = fmt.Sprintf("%s_%s", field, operator)
	}
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

func (qb *QueryBuilder) Where(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, "eq", value)
}

func (qb *QueryBuilder) Or(field, operator string, value interface{}) *QueryBuilder {
	key := field
	if operator != "" && operator != "eq" {
		key = fmt.Sprintf("%s_%s", field, operator)
	}
	filterStr := fmt.Sprintf("[or]%s=%v", key, value)
	qb.orFilters[""] = append(qb.orFilters[""], filterStr)
	return qb
}

func (qb *QueryBuilder) OrGroup(groupName, field, operator string, value interface{}) *QueryBuilder {
	key := field
	if operator != "" && operator != "eq" {
		key = fmt.Sprintf("%s_%s", field, operator)
	}
	filterStr := fmt.Sprintf("[or:%s]%s=%v", groupName, key, value)
	qb.orFilters[groupName] = append(qb.orFilters[groupName], filterStr)
	return qb
}

func (qb *QueryBuilder) MultiFieldOr(fields []string, operator string, value interface{}) *QueryBuilder {
	key := strings.Join(fields, "__or__")
	if operator != "" && operator != "eq" {
		key = fmt.Sprintf("%s_%s", key, operator)
	}
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) Sort(field string) *QueryBuilder {
	qb.sort = field
	return qb
}

func (qb *QueryBuilder) OrderAsc() *QueryBuilder {
	qb.order = "asc"
	return qb
}

func (qb *QueryBuilder) OrderDesc() *QueryBuilder {
	qb.order = "desc"
	return qb
}

func (qb *QueryBuilder) Build() string {
	params := url.Values{}
	
	for key, value := range qb.filters {
		params.Add(key, value)
	}
	
	for _, filters := range qb.orFilters {
		for _, filter := range filters {
			parts := strings.SplitN(filter, "=", 2)
			if len(parts) == 2 {
				params.Add(parts[0], parts[1])
			}
		}
	}
	
	if qb.limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", qb.limit))
	}
	if qb.offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", qb.offset))
	}
	
	if qb.sort != "" {
		params.Add("sort", qb.sort)
		if qb.order != "" {
			params.Add("order", qb.order)
		}
	}
	
	return params.Encode()
}
EOF

# documents.go
cat > cocobase/documents.go << 'EOF'
package cocobase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetDocument(ctx context.Context, collection, docID string) (*Document, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) CreateDocument(ctx context.Context, collection string, data map[string]interface{}) (*Document, error) {
	path := fmt.Sprintf("/collections/documents?collection=%s", collection)
	
	resp, err := c.request(ctx, http.MethodPost, path, data, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) UpdateDocument(ctx context.Context, collection, docID string, data map[string]interface{}) (*Document, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodPatch, path, data, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) DeleteDocument(ctx context.Context, collection, docID string) error {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodDelete, path, nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) ListDocuments(ctx context.Context, collection string, query *QueryBuilder) ([]Document, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	
	if query != nil {
		queryStr := query.Build()
		if queryStr != "" {
			path += "?" + queryStr
		}
	}
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var docs []Document
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return docs, nil
}

func (c *Client) QueryDocuments(ctx context.Context, collection, rawQuery string) ([]Document, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	
	if rawQuery != "" {
		path += "?" + rawQuery
	}
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var docs []Document
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return docs, nil
}
EOF

# auth.go
cat > cocobase/auth.go << 'EOF'
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
EOF

# realtime.go
cat > cocobase/realtime.go << 'EOF'
package cocobase

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func (c *Client) WatchCollection(ctx context.Context, collection string, callback func(Event), name string) (*Connection, error) {
	wsURL := strings.Replace(c.baseURL, "http", "ws", 1)
	wsURL = fmt.Sprintf("%s/realtime/collections/%s", wsURL, collection)
	
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	authMsg := map[string]string{"api_key": c.apiKey}
	if err := conn.WriteJSON(authMsg); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send auth message: %w", err)
	}

	if name == "" {
		name = fmt.Sprintf("watch-%s", collection)
	}

	connection := &Connection{
		conn:   conn,
		name:   name,
		closed: false,
	}

	go func() {
		defer func() {
			connection.mu.Lock()
			connection.closed = true
			connection.mu.Unlock()
		}()

		for {
			var event Event
			err := conn.ReadJSON(&event)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket error: %v\n", err)
				}
				return
			}
			callback(event)
		}
	}()

	return connection, nil
}

func (conn *Connection) Close() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	
	if conn.closed {
		return nil
	}
	
	conn.closed = true
	return conn.conn.Close()
}

func (conn *Connection) IsClosed() bool {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	return conn.closed
}
EOF

# storage/memory.go
cat > storage/memory.go << 'EOF'
package storage

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	value, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	
	return value, nil
}

func (s *MemoryStorage) Set(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.data[key] = value
	return nil
}

func (s *MemoryStorage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.data, key)
	return nil
}
EOF

# storage/file.go
cat > storage/file.go << 'EOF'
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	filepath string
	data     map[string]string
	mu       sync.RWMutex
}

func NewFileStorage(filepath string) (*FileStorage, error) {
	fs := &FileStorage{
		filepath: filepath,
		data:     make(map[string]string),
	}
	
	if err := fs.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	
	return fs, nil
}

func (s *FileStorage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &s.data)
}

func (s *FileStorage) save() error {
	dir := filepath.Dir(s.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(s.filepath, data, 0644)
}

func (s *FileStorage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	value, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	
	return value, nil
}

func (s *FileStorage) Set(key string, value string) error {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
	
	return s.save()
}

func (s *FileStorage) Delete(key string) error {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
	
	return s.save()
}
EOF

# examples/basic/main.go
cat > examples/basic/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Create a document
	doc, err := client.CreateDocument(ctx, "users", map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created document: %s\n", doc.ID)

	// Get a document
	doc, err = client.GetDocument(ctx, "users", doc.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved: %+v\n", doc.Data)

	// List documents with simple query
	query := cocobase.NewQuery().
		Where("age", 30).
		Limit(10)

	docs, err := client.ListDocuments(ctx, "users", query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d users\n", len(docs))
}
EOF

# examples/advanced/main.go
cat > examples/advanced/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Complex query: Active users who are premium OR verified
	query := cocobase.NewQuery().
		Where("status", "active").
		Or("isPremium", "eq", true).
		Or("isVerified", "eq", true).
		Sort("createdAt").
		OrderDesc().
		Limit(50)

	docs, err := client.ListDocuments(ctx, "users", query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d active premium/verified users\n", len(docs))

	// Multi-field search
	query = cocobase.NewQuery().
		MultiFieldOr([]string{"name", "email"}, "contains", "john")

	docs, err = client.ListDocuments(ctx, "users", query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d users matching 'john'\n", len(docs))

	// Multiple OR groups
	query = cocobase.NewQuery().
		OrGroup("tier", "isPremium", "eq", true).
		OrGroup("tier", "isVerified", "eq", true).
		OrGroup("location", "country", "eq", "US").
		OrGroup("location", "country", "eq", "UK")

	docs, err = client.ListDocuments(ctx, "users", query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d users\n", len(docs))
}
EOF

# examples/auth/main.go
cat > examples/auth/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/cocobase-go/cocobase"
	"github.com/yourusername/cocobase-go/storage"
)

func main() {
	// Create client with persistent storage
	store := storage.NewMemoryStorage()
	
	client := cocobase.NewClient(cocobase.Config{
		APIKey:  "your-api-key",
		Storage: store,
	})

	ctx := context.Background()

	// Register a new user
	err := client.Register(ctx, "user@example.com", "password123", map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User registered successfully")

	// Check if authenticated
	if client.IsAuthenticated() {
		fmt.Println("User is authenticated")
	}

	// Get current user
	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current user: %s\n", user.Email)

	// Update user
	newEmail := "newemail@example.com"
	user, err = client.UpdateUser(ctx, map[string]interface{}{
		"phone": "+1234567890",
	}, &newEmail, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated user: %s\n", user.Email)

	// Logout
	if err := client.Logout(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged out")

	// Login
	err = client.Login(ctx, "newemail@example.com", "password123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in successfully")
}
EOF

# examples/realtime/main.go
cat > examples/realtime/main.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Watch collection for changes
	conn, err := client.WatchCollection(ctx, "users", func(event cocobase.Event) {
		fmt.Printf("Event: %s\n", event.Event)
		fmt.Printf("Document ID: %s\n", event.Data.ID)
		fmt.Printf("Data: %+v\n", event.Data.Data)
		fmt.Println("---")
	}, "users-watcher")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Watching for changes... Press Ctrl+C to exit")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nClosing connection...")
}
EOF

# tests/query_test.go
cat > tests/query_test.go << 'EOF'
package tests

import (
	"testing"

	"github.com/yourusername/cocobase-go/cocobase"
)

func TestQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() *cocobase.QueryBuilder
		expected string
	}{
		{
			name: "simple where",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().Where("status", "active")
			},
			expected: "status=active",
		},
		{
			name: "with operator",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().Filter("age", "gte", 18)
			},
			expected: "age_gte=18",
		},
		{
			name: "simple OR",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Or("isPremium", "eq", true).
					Or("isVerified", "eq", true)
			},
			expected: "%5Bor%5DisPremium=true&%5Bor%5DisVerified=true",
		},
		{
			name: "with pagination",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Where("status", "active").
					Limit(50).
					Offset(100)
			},
			expected: "limit=50&offset=100&status=active",
		},
		{
			name: "with sorting",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Sort("createdAt").
					OrderDesc()
			},
			expected: "order=desc&sort=createdAt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder().Build()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMultiFieldOr(t *testing.T) {
	query := cocobase.NewQuery().
		MultiFieldOr([]string{"name", "email"}, "contains", "john")

	result := query.Build()
	expected := "name__or__email_contains=john"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestOrGroups(t *testing.T) {
	query := cocobase.NewQuery().
		OrGroup("tier", "isPremium", "eq", true).
		OrGroup("tier", "isVerified", "eq", true).
		OrGroup("location", "country", "eq", "US").
		OrGroup("location", "country", "eq", "UK")

	result := query.Build()
	
	// Check if all OR group parameters are present
	if !contains(result, "%5Bor%3Atier%5D") {
		t.Error("Missing tier OR group")
	}
	if !contains(result, "%5Bor%3Alocation%5D") {
		t.Error("Missing location OR group")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
EOF

# README.md
cat > README.md << 'EOF'
# Cocobase Go Client

A powerful Go client for the Cocobase Backend as a Service (BaaS).

## Features

- ✅ Full CRUD operations on documents
- ✅ Advanced query filtering with 12+ operators
- ✅ Boolean logic (AND, OR, named OR groups)
- ✅ Multi-field search
- ✅ Authentication (login, register, user management)
- ✅ Real-time updates via WebSocket
- ✅ Pluggable storage for token persistence
- ✅ Thread-safe operations
- ✅ Context support for cancellation and timeouts
- ✅ Comprehensive error handling

## Installation

```bash
go get github.com/yourusername/cocobase-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yourusername/cocobase-go/cocobase"
)

func main() {
    // Initialize client
    client := cocobase.NewClient(cocobase.Config{
        APIKey: "your-api-key",
    })

    ctx := context.Background()

    // Create a document
    doc, err := client.CreateDocument(ctx, "users", map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created: %s\n", doc.ID)
}
```

## Advanced Querying

### Basic Operators

```go
// Equality
query := cocobase.NewQuery().Where("status", "active")

// Comparison operators
query = cocobase.NewQuery().
    Filter("age", "gte", 18).
    Filter("age", "lte", 65)

// String operations
query = cocobase.NewQuery().
    Filter("email", "endswith", "gmail.com")

// IN operator
query = cocobase.NewQuery().
    Filter("role", "in", "admin,moderator,support")

// NULL checks
query = cocobase.NewQuery().
    Filter("deletedAt", "isnull", true)
```

### Boolean Logic

```go
// Simple OR
query := cocobase.NewQuery().
    Where("status", "active").
    Or("isPremium", "eq", true).
    Or("isVerified", "eq", true)

// Named OR groups
query = cocobase.NewQuery().
    OrGroup("tier", "isPremium", "eq", true).
    OrGroup("tier", "isVerified", "eq", true).
    OrGroup("location", "country", "eq", "US").
    OrGroup("location", "country", "eq", "UK")

// Multi-field search
query = cocobase.NewQuery().
    MultiFieldOr([]string{"name", "email"}, "contains", "john")
```

### Pagination & Sorting

```go
query := cocobase.NewQuery().
    Where("status", "active").
    Sort("createdAt").
    OrderDesc().
    Limit(50).
    Offset(100)

docs, err := client.ListDocuments(ctx, "users", query)
```

## Authentication

```go
// Register
err := client.Register(ctx, "user@example.com", "password", map[string]interface{}{
    "firstName": "John",
    "lastName":  "Doe",
})

// Login
err = client.Login(ctx, "user@example.com", "password")

// Get current user
user, err := client.GetCurrentUser(ctx)

// Update user
newEmail := "new@example.com"
user, err = client.UpdateUser(ctx, map[string]interface{}{
    "phone": "+1234567890",
}, &newEmail, nil)

// Logout
err = client.Logout()
```

## Real-time Updates

```go
conn, err := client.WatchCollection(ctx, "users", func(event cocobase.Event) {
    fmt.Printf("Event: %s\n", event.Event)
    fmt.Printf("Document: %+v\n", event.Data)
}, "users-watcher")

if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

## Storage Persistence

```go
import "github.com/yourusername/cocobase-go/storage"

// Memory storage
store := storage.NewMemoryStorage()

// File storage
store, err := storage.NewFileStorage(".cocobase/storage.json")

// Use with client
client := cocobase.NewClient(cocobase.Config{
    APIKey:  "your-api-key",
    Storage: store,
})
```

## Query Operators

| Operator     | Usage                                  |
| ------------ | -------------------------------------- |
| `eq`         | Equals (default)                       |
| `ne`         | Not equals                             |
| `gt`         | Greater than                           |
| `gte`        | Greater than or equal                  |
| `lt`         | Less than                              |
| `lte`        | Less than or equal                     |
| `contains`   | Contains substring (case-insensitive)  |
| `startswith` | Starts with                            |
| `endswith`   | Ends with                              |
| `in`         | In list (comma-separated)              |
| `notin`      | Not in list                            |
| `isnull`     | Is null/not null                       |

## Examples

See the `examples/` directory for complete examples:

- `examples/basic/` - Basic CRUD operations
- `examples/advanced/` - Advanced querying
- `examples/auth/` - Authentication flows
- `examples/realtime/` - WebSocket real-time updates

## Testing

```bash
go test ./tests/...
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
EOF

# .gitignore
cat > .gitignore << 'EOF'
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Storage
.cocobase/
storage.json

# Environment
.env
.env.local
EOF

# Makefile
cat > Makefile << 'EOF'
.PHONY: help build test clean install examples

help:
	@echo "Available commands:"
	@echo "  make install   - Install dependencies"
	@echo "  make build     - Build the project"
	@echo "  make test      - Run tests"
	@echo "  make examples  - Run example programs"
	@echo "  make clean     - Clean build artifacts"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v ./tests/...

examples:
	@echo "Building examples..."
	go build -o bin/basic examples/basic/main.go
	go build -o bin/advanced examples/advanced/main.go
	go build -o bin/auth examples/auth/main.go
	go build -o bin/realtime examples/realtime/main.go

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"
EOF

echo ""
echo "✅ Project structure created successfully!"
echo ""
echo "📋 Next steps:"
echo "   1. Update go.mod with your module path"
echo "   2. Run: go mod download"
echo "   3. Run: go mod tidy"
echo "   4. Update import paths in example files"
echo ""
echo "🚀 Quick commands:"
echo "   make install  - Install dependencies"
echo "   make build    - Build the project"
echo "   make test     - Run tests"
echo ""
echo "📚 Check README.md for usage documentation"