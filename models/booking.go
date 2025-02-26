package models

import "time"

type Booking struct {
	ID          int       `json:"id" db:"id"`
	RideID      int       `json:"ride_id" db:"ride_id" validate:"required"`
	UserID      int       `json:"user_id" db:"user_id" validate:"required"`
	SeatsBooked int       `json:"no_of_seats" db:"no_of_seats" validate:"required,min=1"`
	Status      string    `json:"status" db:"status" validate:"required,oneof=PENDING CONFIRMED CANCELLED"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
