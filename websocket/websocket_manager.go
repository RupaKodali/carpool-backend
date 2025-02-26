package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket connection
type Client struct {
	Conn   *websocket.Conn
	UserID int
	Send   chan []byte
}

// WebSocketManager manages active WebSocket connections
type WebSocketManager struct {
	clients    map[int]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
}

// NewWebSocketManager initializes WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[int]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run starts the WebSocket manager
func (wm *WebSocketManager) Run() {
	for {
		select {
		case client := <-wm.register:
			wm.mu.Lock()
			wm.clients[client.UserID] = client
			wm.mu.Unlock()
			log.Printf("User %d connected\n", client.UserID)

		case client := <-wm.unregister:
			wm.mu.Lock()
			delete(wm.clients, client.UserID)
			wm.mu.Unlock()
			log.Printf("User %d disconnected\n", client.UserID)

		case message := <-wm.broadcast:
			wm.mu.Lock()
			for _, client := range wm.clients {
				client.Send <- message
			}
			wm.mu.Unlock()
		}
	}
}

// SendMessage sends a message to a specific user
func (wm *WebSocketManager) SendMessage(userID int, message []byte) {
	wm.mu.Lock()
	if client, ok := wm.clients[userID]; ok {
		client.Send <- message
	}
	wm.mu.Unlock()
}
