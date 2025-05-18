package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName        string  `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName         string  `json:"last_name" gorm:"type:varchar(100);not null"`
	Username         string  `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Email            string  `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Address          Address `json:"address" gorm:"embedded;embeddedPrefix:user_"`
	Password         string  `json:"password,omitempty" gorm:"type:varchar(255)"`
	Otp              string  `json:"otp,omitempty" gorm:"type:varchar(10)"`
	Phone            string  `gorm:"type:varchar(10);not null"`
	IsDriver         bool    `json:"is_driver" `
	IsEmailVerified  bool    `json:"is_email_verified" `
	IsMobileVerified bool    `json:"is_mobile_verified" `
	GoogleID         string  `json:"google_id,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	AuthProvider     string  `json:"auth_provider" gorm:"type:enum('email','google');default:'email'"`
	LicenseNumber    string  `json:"license_number" gorm:"type:varchar(20)"`
}
