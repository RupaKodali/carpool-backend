package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:UserID;references:ID"`
	RideID      uint
	Ride        Ride   `gorm:"foreignKey:RideID;references:ID"`
	SeatsBooked uint   `gorm:"not null"`
	Status      string `gorm:"type:enum('PENDING','CONFIRMED','CANCELLED');default:PENDING;not null"`
}
