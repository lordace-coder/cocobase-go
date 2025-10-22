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
