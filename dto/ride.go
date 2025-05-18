package dto

import (
	"time"
)

type RideListResponseDTO struct {
	BaseDTO
	DriverID       uint                `json:"driver_id"`
	Driver         UserRideResponseDTO `json:"driver"`
	Origin         LocationDTO         `json:"origin"`
	Destination    LocationDTO         `json:"destination"`
	DepartureAt    time.Time           `json:"departure_at"`
	SeatsAvailable uint                `json:"seats_available"`
	Distance       float64             `json:"distance"`
	DistanceType   string              `json:"distance_type"`
	Duration       string              `json:"duration"`
	Price          float64             `json:"price"`
	Radius         float64             `json:"radius"`
}

type RideResponseDTO struct {
	BaseDTO
	DriverID       uint                `json:"driver_id"`
	Driver         UserRideResponseDTO `json:"driver"`
	Origin         LocationDTO         `json:"origin"`
	Destination    LocationDTO         `json:"destination"`
	DepartureAt    time.Time           `json:"departure_at"`
	SeatsAvailable uint                `json:"seats_available"`
	Route          string              `json:"route,omitempty"`
	Distance       float64             `json:"distance"`
	DistanceType   string              `json:"distance_type"`
	Duration       string              `json:"duration"`
	Price          float64             `json:"price"`
}

type LocationDTO struct {
	FormattedAddress string         `json:"formatted_address"`
	Address          AddressDTO     `json:"address"`
	Coordinates      CoordinatesDTO `json:"coordinates"`
}

type AddressDTO struct {
	Street     string `json:"street"`
	Area       string `json:"area"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type CoordinatesDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
