package models

import "time"

type Message struct {
	ID             int       `json:"id" db:"id"`
	ConversationID int       `json:"conversation_id" db:"conversation_id" validate:"required"`
	SenderID       int       `json:"sender_id" db:"sender_id" validate:"required"`
	ReceiverID     int       `json:"receiver_id" db:"receiver_id" validate:"required"`
	Message        string    `json:"message" db:"message" validate:"required,min=1,max=1000"`
	Status         string    `json:"status" db:"status" validate:"oneof=sent delivered read"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
