package models

import "time"

type Conversation struct {
	ID        int       `json:"id" db:"id"`
	User1ID   int       `json:"user1_id" db:"user1_id" validate:"required"`
	User2ID   int       `json:"user2_id" db:"user2_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
