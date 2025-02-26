package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Group, userController *controllers.UserController, rideController *controllers.RideController, bookingController *controllers.BookingController, messageController *controllers.MessageController, requiredRideController *controllers.RequiredRideController) {
	UserRoutes(e, userController)
	RideRoutes(e, rideController)
	BookingRoutes(e, bookingController)
	MessageRoutes(e, messageController)
	RequiredRideRoutes(e, requiredRideController)
}

func PublicRoutes(e *echo.Echo, userController *controllers.UserController) {
	AuthRoutes(e, userController) // Set up user routes
}
