package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group, userController *controllers.UserController) {
	e.GET("/users/:id", userController.GetUser)       // Get user by ID
	e.PUT("/users/:id", userController.UpdateUser)    // Update user by ID
	e.DELETE("/users/:id", userController.DeleteUser) // Delete user by ID
}

func AuthRoutes(e *echo.Echo, userController *controllers.UserController) {
	e.POST("/users/register", userController.RegisterUser) // Register user
	e.POST("/users/login", userController.LoginUser)       // Login user
}
