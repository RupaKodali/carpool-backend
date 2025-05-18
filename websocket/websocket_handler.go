package websocket

import (
	"carpool-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocketConnection manages WebSocket connections
func HandleWebSocketConnection(wm *WebSocketManager, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	id64, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		log.Println("Invalid User ID")
		conn.Close()
		return
	}

	userID := uint(id64)
	client := &Client{Conn: conn, UserID: userID, Send: make(chan []byte)}
	wm.register <- client

	go client.ReadMessages(db, wm)
	go client.WriteMessages()
}

// ReadMessages listens for incoming messages
func (c *Client) ReadMessages(db *gorm.DB, wm *WebSocketManager) {
	defer func() {
		wm.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		// Parse received message
		var msg map[string]interface{}
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Println("JSON Parse Error:", err)
			continue
		}

		// ✅ Handle Read Receipt
		if msg["type"] == "read" {
			conversationID, ok := msg["conversation_id"].(float64)
			if !ok {
				log.Println("Invalid conversation_id in read receipt")
				continue
			}

			// ✅ Mark messages as "read"
			db.Model(&models.Message{}).
				Where("conversation_id = ? AND receiver_id = ? AND status = 'sent'", int(conversationID), c.UserID).
				Update("status", "read")

			log.Printf("Marked messages as read in conversation %d for user %d", int(conversationID), c.UserID)
			continue
		}

		// ✅ Otherwise, process chat message
		var chatMessage models.Message
		if err := json.Unmarshal(messageBytes, &chatMessage); err != nil {
			log.Println("Message Parse Error:", err)
			continue
		}

		// Save message to database
		chatMessage.CreatedAt = time.Now()
		chatMessage.Status = "sent"
		if err := db.Create(&chatMessage).Error; err != nil {
			log.Println("DB Insert Error:", err)
			continue
		}

		// Send message to the receiver if they are online
		wm.SendMessage(chatMessage.ReceiverID, messageBytes)
	}
}

// WriteMessages sends messages to the client
func (c *Client) WriteMessages() {
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write Error:", err)
			break
		}
	}
}
