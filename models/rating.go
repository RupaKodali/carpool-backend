package models

import "time"

type Rating struct {
	ID        int       `json:"id" db:"id"`
	RaterID   int       `json:"rater_id" db:"rater_id"`
	RateeID   int       `json:"ratee_id" db:"ratee_id"`
	RideID    int       `json:"ride_id" db:"ride_id"`
	Rating    int       `json:"rating" db:"rating"`           // Range: 1 to 5
	Review    *string   `json:"review,omitempty" db:"review"` // Pointer for nullable field
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
