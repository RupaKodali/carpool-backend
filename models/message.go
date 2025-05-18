package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ConversationID uint
	Conversation   Conversation `gorm:"foreignKey:ConversationID;references:ID"`
	SenderID       uint
	Sender         User `gorm:"foreignKey:SenderID;references:ID"`
	ReceiverID     uint
	Receiver       User   `gorm:"foreignKey:ReceiverID;references:ID"`
	Message        string `gorm:"type:text;not null"`
	Status         string `gorm:"type:enum('sent','delivered','read');default:sent"`
}
