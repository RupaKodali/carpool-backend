package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

func RequiredRideRoutes(e *echo.Group, requiredRideController *controllers.RequiredRideController) {
	e.POST("/required-rides", requiredRideController.CreateRequiredRide)       // Create a new required ride
	e.GET("/required-rides", requiredRideController.ListRequiredRides)         // List all required rides
	e.DELETE("/required-rides/:id", requiredRideController.DeleteRequiredRide) // Delete a required ride by ID
}
