package controllers

import (
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookingController struct {
	BookingService services.BookingService
}

// NewBookingController creates a new BookingController with the given BookingService
func NewBookingController(BookingService services.BookingService) *BookingController {
	return &BookingController{BookingService: BookingService}
}

// CreateBooking handles POST /bookings
func (h *BookingController) CreateBooking(c echo.Context) error {

	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var booking models.Booking
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	booking.UserID = loggedInUserID

	booking.Status = "CONFIRMED"

	err = h.BookingService.CreateBooking(&booking)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Booking created successfully"})
}

// GetBooking handles GET /bookings/:id
func (h *BookingController) GetBooking(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid booking ID"})
	}

	id := uint(id64)

	booking, err := h.BookingService.GetBookingByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, booking)
}

// DeleteBooking handles DELETE /bookings/:id
func (h *BookingController) DeleteBooking(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get booking ID from request param
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid booking ID"})
	}

	id := uint(id64)
	// Check if the logged-in user is the owner of the booking
	booking, err := h.BookingService.GetBookingByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Booking not found"})
	}
	if booking.UserID != loggedInUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to delete this booking"})
	}

	err = h.BookingService.DeleteBooking(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Booking deleted successfully"})
}

// ListBookings handles fetching bookings dynamically based on query parameters
func (h *BookingController) ListBookings(c echo.Context) error {
	params := services.ParseQueryParams(c)

	// Check if ride_id is present in query params
	if rideID := c.QueryParam("ride_id"); rideID != "" {
		id, err := strconv.Atoi(rideID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ride ID"})
		}
		params.Filters["ride_id"] = id
	} else {
		// If no ride_id is provided, fetch only bookings for the logged-in user
		userID, err := utils.GetUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
		}
		params.Filters["user_id"] = userID
	}

	// Call the service function
	bookings, err := h.BookingService.ListBookings(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bookings)
}
