package services

import (
	"carpool-backend/models"
	"errors"

	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(booking *models.Booking) error
	GetBookingByID(id int) (*models.Booking, error)
	DeleteBooking(id int) error
	ListBookings(params QueryParams) (*PaginatedResponse, error)
}

type bookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) BookingService {
	return &bookingService{db: db}
}

func (s *bookingService) CreateBooking(booking *models.Booking) error {
	var ride models.Ride

	// Check if the ride exists and has available seats
	if err := s.db.First(&ride, booking.RideID).Error; err != nil {
		return errors.New("ride not found")
	}
	if ride.SeatsAvailable < booking.SeatsBooked {
		return errors.New("not enough seats available")
	}

	// Begin transaction
	tx := s.db.Begin()

	// Insert the booking
	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create booking")
	}

	// Update available seats
	if err := tx.Model(&ride).Update("seats_available", ride.SeatsAvailable-booking.SeatsBooked).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update available seats")
	}

	tx.Commit()
	return nil
}

func (s *bookingService) GetBookingByID(id int) (*models.Booking, error) {
	var booking models.Booking
	if err := s.db.First(&booking, id).Error; err != nil {
		return nil, errors.New("booking not found")
	}
	return &booking, nil
}

func (s *bookingService) DeleteBooking(id int) error {
	var booking models.Booking

	// Get booking details
	if err := s.db.First(&booking, id).Error; err != nil {
		return errors.New("booking not found")
	}

	// Start transaction
	tx := s.db.Begin()

	// Delete the booking
	if err := tx.Delete(&booking).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete booking")
	}

	// Update available seats
	if err := tx.Model(&models.Ride{}).Where("id = ?", booking.RideID).
		Update("seats_available", gorm.Expr("seats_available + ?", booking.SeatsBooked)).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update available seats")
	}

	tx.Commit()
	return nil
}

// ListBookings fetches bookings dynamically based on filters, pagination, search, and sorting
func (s *bookingService) ListBookings(params QueryParams) (*PaginatedResponse, error) {
	var bookings []models.Booking
	searchableFields := []string{}

	return ListEntities(s.db, &bookings, params, searchableFields)
}
