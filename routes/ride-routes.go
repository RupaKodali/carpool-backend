package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

func RideRoutes(e *echo.Group, rideController *controllers.RideController) {
	e.POST("/rides", rideController.CreateRide)       // Create a new ride
	e.GET("/rides/:id", rideController.GetRide)       // Get a ride by ID
	e.PUT("/rides/:id", rideController.UpdateRide)    // Update a ride by ID
	e.DELETE("/rides/:id", rideController.DeleteRide) // Delete a ride by ID
	e.GET("/rides", rideController.ListRides)         // List all rides
	e.POST("/rides/match", rideController.MatchRides) // Add ride matching endpoint

}
