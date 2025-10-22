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
