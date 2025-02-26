package websocket

import (
	"carpool-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		log.Println("Invalid User ID")
		conn.Close()
		return
	}

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

		// Parse message JSON into the Message struct
		var msg models.Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Println("JSON Parse Error:", err)
			continue
		}

		// Store message using GORM
		newMessage := models.Message{
			ConversationID: msg.ConversationID,
			SenderID:       msg.SenderID,
			ReceiverID:     msg.ReceiverID,
			Message:        msg.Message,
			Status:         "sent",
		}

		// Insert message into the database
		if err := db.Create(&newMessage).Error; err != nil {
			log.Println("DB Insert Error:", err)
			continue
		}

		// Send message to the recipient via WebSocket
		wm.SendMessage(msg.ReceiverID, messageBytes)
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
