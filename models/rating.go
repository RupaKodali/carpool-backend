package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	RaterID uint
	RateeID uint
	RideID  uint
	Rater   User    `gorm:"foreignKey:RaterID;references:ID"`
	Ratee   User    `gorm:"foreignKey:RateeID;references:ID"`
	Ride    Ride    `gorm:"foreignKey:RideID;references:ID"`
	Rating  uint    `gorm:"not null"`
	Review  *string `gorm:"type:text"`
}
