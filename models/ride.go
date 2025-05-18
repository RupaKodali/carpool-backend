package models

import (
	"time"

	"gorm.io/gorm"
)

type Ride struct {
	gorm.Model
	DriverID       uint
	Driver         User      `gorm:"foreignKey:DriverID;references:ID"`
	Origin         Location  `gorm:"embedded;embeddedPrefix:origin_"`
	Destination    Location  `gorm:"embedded;embeddedPrefix:destination_"`
	DepartureAt    time.Time `json:"departure_at" gorm:"not null"`
	SeatsAvailable uint      `json:"seats_available" gorm:"not null"`
	Route          string    `json:"route" gorm:"type:longtext"`
	Distance       float64   `json:"distance" `
	DistanceType   string    `json:"distance_type" gorm:"type:longtext"`
	Duration       string    `json:"duration" gorm:"type:longtext"`
	Price          float64   `json:"price" `
}

type Address struct {
	Street     string `json:"street" gorm:"column:street;type:varchar(255)"`
	Area       string `json:"area" gorm:"column:area;type:varchar(255)"`
	City       string `json:"city" gorm:"column:city;type:varchar(255)"`
	State      string `json:"state" gorm:"column:state;type:varchar(100)"`
	Country    string `json:"country" gorm:"column:country;type:varchar(100)"`
	PostalCode string `json:"postalCode" gorm:"column:postal_code;type:varchar(20)"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude" gorm:"column:latitude"`
	Longitude float64 `json:"longitude" gorm:"column:longitude"`
}

type Location struct {
	FormattedAddress string      `json:"formatted_address" gorm:"column:formatted_address;type:varchar(255)"`
	Address          Address     `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	Coordinates      Coordinates `json:"coordinates" gorm:"embedded;embeddedPrefix:coordinates_"`
}
