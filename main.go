package main

import (
	"log"
	"os"

	"carpool-backend/controllers"
	"carpool-backend/database"
	"carpool-backend/routes"
	"carpool-backend/services"
	"carpool-backend/utils"
	"carpool-backend/websocket"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create a new Echo instance
	e := echo.New()

	// Register custom validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// WebSocket Setup
	wm := websocket.NewWebSocketManager()
	go wm.Run()

	// ✅ Fixed: Pass `db` along with `wm`
	e.GET("/ws", func(c echo.Context) error {
		websocket.HandleWebSocketConnection(wm, db, c.Response().Writer, c.Request())
		return nil
	})

	// Initialize services
	userService := services.NewUserService(db)
	rideService := services.NewRideService(db)
	bookingService := services.NewBookingService(db)
	messageService := services.NewMessageService(db)
	requiredRideService := services.NewRequiredRideService(db)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	rideController := controllers.NewRideController(rideService)
	bookingController := controllers.NewBookingController(bookingService)
	messageController := controllers.NewMessageController(messageService, wm)
	requiredRideController := controllers.NewRequiredRideController(requiredRideService)

	// Public routes
	routes.PublicRoutes(e, userController)

	// JWT Middleware for protected routes
	jwtSecret := os.Getenv("JWT_SECRET")
	authGroup := e.Group("/auth")
	authGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
	}))

	// Set up protected routes
	routes.SetupRoutes(authGroup, userController, rideController, bookingController, messageController, requiredRideController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
