package dto

import (
	"time"
)

type RequiredRideDTO struct {
	BaseDTO
	UserID      uint
	User        UserRideResponseDTO `json:"user"`
	Origin      LocationDTO         `json:"origin"`
	Destination LocationDTO         `json:"destination"`
	DepartureAt time.Time           `json:"departure_at"`
	Radius      float64             `json:"radius"`
}
