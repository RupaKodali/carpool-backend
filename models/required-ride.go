package models

import (
	"time"

	"gorm.io/gorm"
)

type RequiredRide struct {
	gorm.Model
	UserID      uint
	User        User      `gorm:"foreignKey:UserID;references:ID"`
	Origin      Location  `gorm:"embedded;embeddedPrefix:origin_"`
	Destination Location  `gorm:"embedded;embeddedPrefix:destination_"`
	DepartureAt time.Time `gorm:"not null"`
}
