package controllers

import (
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type RideController struct {
	RideService services.RideService
}

// NewRideController creates a new RideController with the given UserService
func NewRideController(rideService services.RideService) *RideController {
	return &RideController{RideService: rideService}
}

// CreateRide handles POST /rides
func (h *RideController) CreateRide(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	isDriver, _ := utils.IsDriverFromToken(c)
	if !isDriver {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Not a driver"})
	}
	var ride models.Ride
	if err := c.Bind(&ride); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	ride.DriverID = loggedInUserID
	err = h.RideService.CreateRide(&ride)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Ride created successfully"})
}

// GetRide handles GET /rides/:id
func (h *RideController) GetRide(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ride ID"})
	}

	ride, err := h.RideService.GetRideByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ride)
}

// UpdateRide handles PUT /rides/:id
func (h *RideController) UpdateRide(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get ride ID from request param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ride ID"})
	}

	// Check if the logged-in user is the owner of the ride
	ride, err := h.RideService.GetRideByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ride not found"})
	}
	if ride.DriverID != loggedInUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to update this ride"})
	}

	// Proceed with updating the ride
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Ensure ID is not modified
	updates["id"] = id

	err = h.RideService.UpdateRide(ride, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ride updated successfully"})
}

// DeleteRide handles DELETE /rides/:id
func (h *RideController) DeleteRide(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get ride ID from request param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ride ID"})
	}

	// Check if the logged-in user is the owner of the ride
	ride, err := h.RideService.GetRideByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ride not found"})
	}
	if ride.DriverID != loggedInUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to update this ride"})
	}

	err = h.RideService.DeleteRide(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ride deleted successfully"})
}

// ListRides handles GET /rides
func (h *RideController) ListRides(c echo.Context) error {
	params := services.ParseQueryParams(c)

	// Parse date filters
	if departureAt := c.QueryParam("departure_at"); departureAt != "" {

		// If user provides `departure_at`, parse it
		if date, err := time.Parse(time.RFC3339, departureAt); err == nil {
			params.Filters["departure_at"] = map[string]interface{}{"from": date}
		}
	} else {
		// If `departure_at` is NOT provided, default to today
		today := time.Now().Format("2006-01-02") // Format as YYYY-MM-DD
		params.Filters["departure_at"] = map[string]interface{}{"from": today}
	}

	// Fetch rides with dynamic filters
	response, err := h.RideService.ListRides(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// MatchRides handles POST /auth/rides/match
func (h *RideController) MatchRides(c echo.Context) error {
	// MatchRideRequest represents the structure for ride matching criteria
	var request struct {
		OriginLat      float64    `json:"origin_lat" validate:"required"`
		OriginLng      float64    `json:"origin_lng" validate:"required"`
		DestinationLat float64    `json:"destination_lat" validate:"required"`
		DestinationLng float64    `json:"destination_lng" validate:"required"`
		DepartureAt    *time.Time `json:"departure_at"`
		Radius         *float64   `json:"radius"` // Optional, defaults to 0.5 km if not provided

	}

	// Bind and validate the request
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	radius := 0.5
	if request.Radius != nil && *request.Radius > 0 {
		radius = *request.Radius
	}

	// Initialize query parameters
	params := services.ParseQueryParams(c)

	// Time-based filtering (±2 hours from given departure time)
	const maxTimeDiffHours = 24.0 // Maximum acceptable time difference in hours
	if request.DepartureAt != nil {
		// startTime := request.DepartureAt.Add(-time.Hour * time.Duration(maxTimeDiffHours))
		// endTime := request.DepartureAt.Add(time.Hour * time.Duration(maxTimeDiffHours))
		// params.Filters["departure_at"] = map[string]interface{}{
		// 	"from": startTime,
		// 	"to":   endTime,
		// }
		now := time.Now().In(request.DepartureAt.Location())
		requestDate := request.DepartureAt

		isSameDay := now.Year() == requestDate.Year() &&
			now.Month() == requestDate.Month() &&
			now.Day() == requestDate.Day()

		var fromTime, toTime time.Time

		if isSameDay {
			fromTime = now
			toTime = time.Date(
				requestDate.Year(), requestDate.Month(), requestDate.Day(),
				23, 59, 59, int(time.Second-time.Nanosecond), requestDate.Location(),
			)
		} else {
			fromTime = time.Date(
				requestDate.Year(), requestDate.Month(), requestDate.Day(),
				0, 0, 0, 0, requestDate.Location(),
			)
			toTime = fromTime.Add(24 * time.Hour).Add(-time.Nanosecond)
		}

		params.Filters["departure_at"] = map[string]interface{}{
			"from": fromTime,
			"to":   toTime,
		}
	}

	// Fetch available rides
	result, err := h.RideService.ListRides(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch available rides"})
	}

	// Extract rides from paginated response
	var availableRides []models.Ride
	if ridesPtr, ok := result.Data.(*[]models.Ride); ok {
		availableRides = *ridesPtr // Dereference the pointer to get the slice
	} else {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Invalid data format"})
	}

	// Match rides using the  geolocation-based matching
	matchingRides, err := services.MatchRides(
		request.OriginLat, request.OriginLng,
		request.DestinationLat, request.DestinationLng,
		radius,
		availableRides,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "No matching rides found"})
	}

	// Update response data with matched rides
	result.Data = matchingRides

	return c.JSON(http.StatusOK, result)
}
