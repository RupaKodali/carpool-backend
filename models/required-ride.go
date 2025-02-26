package models

import "time"

type RequiredRide struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id" validate:"required"`
	Origin      string    `json:"origin" db:"origin" validate:"required,min=3,max=255"`
	Destination string    `json:"destination" db:"destination" validate:"required,min=3,max=255"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
