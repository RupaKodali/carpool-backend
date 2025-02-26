package models

import "time"

type User struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email         string    `json:"email" db:"email" validate:"required,email"`
	Password      string    `json:"password,omitempty" db:"password" validate:"required,min=8"`
	Phone         string    `json:"phone" db:"phone" validate:"required,len=10,numeric"`
	IsDriver      bool      `json:"is_driver" db:"is_driver"`
	LicenseNumber *string   `json:"license_number,omitempty" db:"license_number" validate:"omitempty,min=5,max=20"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
