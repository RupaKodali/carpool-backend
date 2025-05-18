package dto

import (
	"time"
)

type UserRideResponseDTO struct {
	BaseDTO
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type BaseDTO struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserLoginResponse struct {
	BaseDTO
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	IsEmailVerified  bool
	IsMobileVerified bool
	IsDriver         bool
}

type TokenStruct struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserResponseDTO struct {
	BaseDTO
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	IsEmailVerified  bool
	IsMobileVerified bool
	IsDriver         bool
}
