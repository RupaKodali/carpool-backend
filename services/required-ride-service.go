package services

import (
	"carpool-backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RequiredRideService interface {
	CreateRequiredRide(ride *models.RequiredRide) error
	ListRequiredRides() ([]models.RequiredRide, error)
	DeleteRequiredRide(id int) error
	GetRequiredRides(id int) (*models.RequiredRide, error)
}

type requiredRideService struct {
	db *gorm.DB
}

func NewRequiredRideService(db *gorm.DB) RequiredRideService {
	return &requiredRideService{db: db}
}

func (s *requiredRideService) CreateRequiredRide(ride *models.RequiredRide) error {
	ride.CreatedAt = time.Now()

	// Insert into database using GORM
	if err := s.db.Create(&ride).Error; err != nil {
		return errors.New("failed to create required ride")
	}
	return nil
}

func (s *requiredRideService) ListRequiredRides() ([]models.RequiredRide, error) {
	var rides []models.RequiredRide

	// Fetch required rides in descending order of creation
	if err := s.db.Order("created_at DESC").Find(&rides).Error; err != nil {
		return nil, errors.New("failed to list required rides")
	}
	return rides, nil
}

func (s *requiredRideService) GetRequiredRides(id int) (*models.RequiredRide, error) {
	var ride models.RequiredRide

	// Fetch ride by ID
	if err := s.db.First(&ride, id).Error; err != nil {
		return nil, errors.New("ride not found")
	}
	return &ride, nil
}

func (s *requiredRideService) DeleteRequiredRide(id int) error {
	// Delete the required ride
	if err := s.db.Delete(&models.RequiredRide{}, id).Error; err != nil {
		return errors.New("failed to delete required ride")
	}
	return nil
}
