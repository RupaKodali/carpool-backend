package dto

import "time"

type RideListResponseDTO struct {
	ID       int                 `json:"id" db:"id"`
	DriverID int                 `json:"driver_id"`
	Driver   UserRideResponseDTO `json:"driver" `
	Origin   string              `json:"origin"`
	// OriginLat      float64             `json:"origin_lat" `
	// OriginLng      float64             `json:"origin_lng" `
	Destination string `json:"destination" `
	// DestinationLat float64             `json:"destination_lat" `
	// DestinationLng float64             `json:"destination_lng" `
	DepartureAt    time.Time `json:"departure_at" `
	SeatsAvailable int       `json:"seats_available"`
	// Route          string              `json:"route,omitempty" `
	Distance     float64   `json:"distance" `
	DistanceType string    `json:"distance_type" `
	Duration     string    `json:"duration" `
	Price        float64   `json:"price" `
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RideResponseDTO struct {
	ID             int                 `json:"id" db:"id"`
	DriverID       int                 `json:"driver_id"`
	Driver         UserRideResponseDTO `json:"driver" `
	Origin         string              `json:"origin"`
	OriginLat      float64             `json:"origin_lat" `
	OriginLng      float64             `json:"origin_lng" `
	Destination    string              `json:"destination" `
	DestinationLat float64             `json:"destination_lat" `
	DestinationLng float64             `json:"destination_lng" `
	DepartureAt    time.Time           `json:"departure_at" `
	SeatsAvailable int                 `json:"seats_available"`
	Route          string              `json:"route,omitempty" `
	Distance       float64             `json:"distance" `
	DistanceType   string              `json:"distance_type" `
	Duration       string              `json:"duration" `
	Price          float64             `json:"price" `
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}
