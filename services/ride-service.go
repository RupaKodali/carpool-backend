package services

import (
	"carpool-backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RideService interface {
	CreateRide(ride *models.Ride) error
	GetRideByID(id uint, preloads ...string) (*models.Ride, error)
	UpdateRide(existingRide *models.Ride, updates map[string]interface{}) error
	DeleteRide(id uint) error
	ListRides(params QueryParams) (*PaginatedResponse, error)
}

type rideService struct {
	db *gorm.DB
}

// NewRideService creates a new RideService instance
func NewRideService(db *gorm.DB) RideService {
	return &rideService{db: db}
}

// CreateRide inserts a new ride into the database
func (s *rideService) CreateRide(ride *models.Ride) error {
	ride.CreatedAt = time.Now()
	ride.UpdatedAt = time.Now()

	if err := s.db.Create(&ride).Error; err != nil {
		return errors.New("failed to create ride")
	}
	return nil
}

// GetRideByID retrieves a ride by its ID
func (s *rideService) GetRideByID(id uint, preloads ...string) (*models.Ride, error) {
	var ride models.Ride

	db := s.db

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	if err := db.First(&ride, id).Error; err != nil {
		return nil, errors.New("ride not found")
	}
	if err := s.db.First(&ride, id).Error; err != nil {
		return nil, errors.New("ride not found")
	}
	return &ride, nil
}

func (s *rideService) UpdateRide(existingRide *models.Ride, updates map[string]interface{}) error {
	// Preserve `created_at`
	updates["created_at"] = existingRide.CreatedAt

	// Ensure `updated_at` is always set
	updates["updated_at"] = time.Now()

	// Perform the update using GORM's `Updates()`
	if err := s.db.Model(existingRide).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DeleteRide deletes a ride by its ID
func (s *rideService) DeleteRide(id uint) error {
	if err := s.db.Delete(&models.Ride{}, id).Error; err != nil {
		return errors.New("failed to delete ride")
	}
	return nil
}

func (s *rideService) ListRides(params QueryParams) (*PaginatedResponse, error) {
	var rides []models.Ride
	searchableFields := []string{"origin", "destination"}

	return ListEntities(s.db, &rides, params, searchableFields)
}
