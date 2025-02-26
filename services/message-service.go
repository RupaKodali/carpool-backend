package services

import (
	"carpool-backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// MessageService defines methods for handling chat messages
type MessageService interface {
	SendMessage(message *models.Message) error
	GetMessageHistory(userID, otherUserID int, params QueryParams) (*PaginatedResponse, error)
	MarkMessagesAsRead(userID, conversationID int) error
	GetConversations(userID int) ([]map[string]interface{}, error)
}

type chatService struct {
	db *gorm.DB
}

// NewMessageService initializes a new MessageService
func NewMessageService(db *gorm.DB) MessageService {
	return &chatService{db: db}
}

func (s *chatService) SendMessage(message *models.Message) error {
	message.CreatedAt = time.Now()

	// Check if conversation already exists
	var conversation models.Conversation
	err := s.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		message.SenderID, message.ReceiverID, message.ReceiverID, message.SenderID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		// Create a new conversation
		newConversation := models.Conversation{
			User1ID:   message.SenderID,
			User2ID:   message.ReceiverID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.db.Create(&newConversation).Error; err != nil {
			return errors.New("failed to create conversation")
		}
		conversation = newConversation
	}

	// Assign conversation_id to the message
	message.ConversationID = conversation.ID

	// Insert the message using GORM
	if err := s.db.Create(&message).Error; err != nil {
		return errors.New("failed to send message")
	}

	return nil
}

// GetMessageHistory retrieves messages between two users
func (s *chatService) GetMessageHistory(userID, otherUserID int, params QueryParams) (*PaginatedResponse, error) {
	var messages []models.Message
	searchableFields := []string{"message"}

	// Find the conversation between the two users
	var conversation models.Conversation
	err := s.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, otherUserID, otherUserID, userID).
		First(&conversation).Error

	if err != nil {
		return nil, errors.New("conversation not found")
	}

	// Add conversation_id filter
	params.Filters["conversation_id"] = conversation.ID

	// Use ListEntities to fetch paginated messages
	return ListEntities(s.db, &messages, params, searchableFields)
}

// MarkMessagesAsRead updates the status of unread messages to "read"
func (s *chatService) MarkMessagesAsRead(userID, conversationID int) error {
	// Update unread messages for the given conversation
	err := s.db.Model(&models.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND status != 'read'", conversationID, userID).
		Update("status", "read").Error

	if err != nil {
		return errors.New("failed to mark messages as read")
	}
	return nil
}

func (s *chatService) GetConversations(userID int) ([]map[string]interface{}, error) {
	var conversations []models.Conversation
	var result []map[string]interface{}

	// Fetch all conversations where user is involved
	err := s.db.Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("updated_at DESC").Find(&conversations).Error
	if err != nil {
		return nil, errors.New("failed to retrieve conversations")
	}

	for _, conv := range conversations {
		// Determine chat partner
		chatPartnerID := conv.User1ID
		if userID == conv.User1ID {
			chatPartnerID = conv.User2ID
		}

		// Fetch chat partner details
		var chatPartner models.User
		err := s.db.Where("id = ?", chatPartnerID).First(&chatPartner).Error
		if err != nil {
			return nil, errors.New("failed to retrieve chat partner details")
		}

		// Fetch last message in conversation
		var lastMessage models.Message
		err = s.db.Where("conversation_id = ?", conv.ID).
			Order("created_at DESC").First(&lastMessage).Error
		if err != nil {
			lastMessage.Message = "No messages yet"
			lastMessage.CreatedAt = conv.CreatedAt
		}

		// Count unread messages
		var unreadCount int64
		s.db.Model(&models.Message{}).
			Where("conversation_id = ? AND receiver_id = ? AND status = 'sent'", conv.ID, userID).
			Count(&unreadCount)

		// Construct response
		conversationData := map[string]interface{}{
			"conversation_id": conv.ID,
			"chat_partner": map[string]interface{}{
				"id":    chatPartner.ID,
				"name":  chatPartner.Name,
				"email": chatPartner.Email,
			},
			"last_message": map[string]interface{}{
				"text":      lastMessage.Message,
				"timestamp": lastMessage.CreatedAt,
			},
			"unread_count": unreadCount,
		}

		result = append(result, conversationData)
	}

	return result, nil
}
