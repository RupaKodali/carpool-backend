package routes

import (
	"carpool-backend/controllers"

	"github.com/labstack/echo/v4"
)

// MessageRoutes defines routes for chat messaging
func MessageRoutes(e *echo.Group, chatController *controllers.MessageController) {
	e.POST("/messages", chatController.SendMessage)               // Send a message
	e.GET("/messages/:user_id", chatController.GetMessageHistory) // Get chat history with a specific user
	e.PUT("/messages/read", chatController.MarkMessagesAsRead)    // Mark messages as read
}
