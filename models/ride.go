package models

import "time"

type Ride struct {
	ID             int       `json:"id" db:"id"`
	DriverID       int       `json:"driver_id" db:"driver_id" validate:"required"`
	Origin         string    `json:"origin" db:"origin" validate:"required,min=3,max=255"`
	OriginLat      float64   `json:"origin_lat" db:"origin_lat" validate:"required"`
	OriginLng      float64   `json:"origin_lng" db:"origin_lng" validate:"required"`
	Destination    string    `json:"destination" db:"destination" validate:"required,min=3,max=255"`
	DestinationLat float64   `json:"destination_lat" db:"destination_lat" validate:"required"`
	DestinationLng float64   `json:"destination_lng" db:"destination_lng" validate:"required"`
	DepartureAt    time.Time `json:"departure_at" db:"departure_at" validate:"required"`
	SeatsAvailable int       `json:"seats_available" db:"seats_available" validate:"required,min=1,max=7"`
	Route          string    `json:"route,omitempty" db:"route" validate:"omitempty"`
	Distance       float64   `json:"distance" db:"distance" validate:"omitempty"`
	DistanceType   string    `json:"distance_type" db:"distance_type" validate:"omitempty"`
	Duration       string    `json:"duration" db:"duration" validate:"omitempty"`
	Price          float64   `json:"price" db:"price" validate:"omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
