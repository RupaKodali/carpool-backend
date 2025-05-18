package models

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	User1ID uint
	User2ID uint
	User1   User `gorm:"foreignKey:User1ID;references:ID"`
	User2   User `gorm:"foreignKey:User2ID;references:ID"`
}
