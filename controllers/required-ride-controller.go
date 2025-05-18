package controllers

import (
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RequiredRideController struct {
	RequiredRideService services.RequiredRideService
}

// NewRequiredRideController creates a new RequiredRideController with the given RequiredRideService
func NewRequiredRideController(RequiredRideService services.RequiredRideService) *RequiredRideController {
	return &RequiredRideController{RequiredRideService: RequiredRideService}
}

// CreateRequiredRide handles POST /required-rides
func (h *RequiredRideController) CreateRequiredRide(c echo.Context) error {
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var ride models.RequiredRide
	if err := c.Bind(&ride); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	ride.ID = loggedInUserID
	err = h.RequiredRideService.CreateRequiredRide(&ride)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Required ride created successfully"})
}

// ListRequiredRides handles GET /required-rides
func (h *RequiredRideController) ListRequiredRides(c echo.Context) error {
	rides, err := h.RequiredRideService.ListRequiredRides()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rides)
}

// DeleteRequiredRide handles DELETE /required-rides/:id
func (h *RequiredRideController) DeleteRequiredRide(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get ride ID from request param
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid booking ID"})
	}

	id := uint(id64)
	// Check if the logged-in user is the owner of the ride
	ride, err := h.RequiredRideService.GetRequiredRides(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Required ride not found"})
	}
	if ride.ID != loggedInUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to delete this required ride"})
	}

	err = h.RequiredRideService.DeleteRequiredRide(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Required ride deleted successfully"})
}
