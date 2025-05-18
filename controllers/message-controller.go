package controllers

import (
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"carpool-backend/websocket"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// MessageController handles chat messaging
type MessageController struct {
	MessageService   services.MessageService
	WebSocketManager *websocket.WebSocketManager
}

// NewMessageController initializes MessageController
func NewMessageController(MessageService services.MessageService, wm *websocket.WebSocketManager) *MessageController {
	return &MessageController{
		MessageService:   MessageService,
		WebSocketManager: wm,
	}
}

// SendMessage handles sending a message (POST /messages)
func (h *MessageController) SendMessage(c echo.Context) error {
	var message models.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Retrieve sender ID from JWT claims
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	message.SenderID = userID

	// Store message in DB
	err = h.MessageService.SendMessage(&message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Send message via WebSocket if receiver is online
	messageJSON, _ := json.Marshal(message)
	h.WebSocketManager.SendMessage(message.ReceiverID, messageJSON)

	return c.JSON(http.StatusCreated, echo.Map{"message": "Message sent successfully"})
}

// GetMessageHistory handles fetching chat history (GET /messages/:user_id)
func (h *MessageController) GetMessageHistory(c echo.Context) error {
	// Get current logged-in user ID from JWT
	currentUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get chat partner's user ID from URL
	UserId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}
	chatUserID := uint(UserId)

	// Parse query parameters for pagination & search
	params := services.ParseQueryParams(c)

	// Fetch messages with filters
	response, err := h.MessageService.GetMessageHistory(currentUserID, chatUserID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// MarkMessagesAsRead handles marking messages as read (PUT /messages/read)
func (h *MessageController) MarkMessagesAsRead(c echo.Context) error {
	var request struct {
		ConversationID uint `json:"conversation_id"`
	}

	UserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	err = h.MessageService.MarkMessagesAsRead(UserID, request.ConversationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update message status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Messages marked as read"})
}

func (h *MessageController) GetConversations(c echo.Context) error {
	// Get current logged-in user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Fetch user's conversations
	conversations, err := h.MessageService.GetConversations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch conversations"})
	}

	return c.JSON(http.StatusOK, conversations)
}
