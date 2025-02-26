package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

func BookingRoutes(e *echo.Group, bookingController *controllers.BookingController) {
	e.POST("/bookings", bookingController.CreateBooking)       // Create a new booking
	e.GET("/bookings/:id", bookingController.GetBooking)       // Get booking by ID
	e.DELETE("/bookings/:id", bookingController.DeleteBooking) // Delete a booking by ID
	e.GET("/bookings", bookingController.ListBookings)         // List all bookings for a specific ride
}
