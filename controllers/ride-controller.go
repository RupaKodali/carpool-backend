package controllers

import (
	"carpool-backend/dto"
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

	ride, err := h.RideService.GetRideByID(id, "Driver")
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	dtoRide, err := dto.MapToDTOs[models.Ride, dto.RideResponseDTO]([]models.Ride{*ride})
	if err != nil || len(dtoRide) == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to map ride to DTO"})
	}

	return c.JSON(http.StatusOK, dtoRide[0])
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

	params.Preloads = append(params.Preloads, "Driver")

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
		FromDateTime   *time.Time `json:"from_datetime" validate:"required"`
		ToDateTime     *time.Time `json:"to_datetime" validate:"required"`
		// DepartureAt    *time.Time `json:"departure_at"`
		Radius *float64 `json:"radius"` // Optional, defaults to 0.5 km if not provided

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

	now := time.Now()
	var from, to time.Time

	if request.FromDateTime != nil {
		fromDate := *request.FromDateTime

		// If only date was given (e.g., time is 00:00:00), or it's today and time has passed
		if fromDate.Hour() == 0 && fromDate.Minute() == 0 && fromDate.Second() == 0 {
			if fromDate.Year() == now.Year() && fromDate.YearDay() == now.YearDay() {
				from = now // current time onwards today
			} else {
				// start of the given day
				from = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
			}
		} else {
			if fromDate.Before(now) {
				from = now
			} else {
				from = fromDate
			}
		}
	} else {
		from = now
	}

	if request.ToDateTime != nil {
		toDate := *request.ToDateTime
		// If no time part was given, assume end of that day
		if toDate.Hour() == 0 && toDate.Minute() == 0 && toDate.Second() == 0 {
			to = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), toDate.Location())
		} else {
			to = toDate
		}
	} else {
		// Default to end of from_date if not given
		to = time.Date(from.Year(), from.Month(), from.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), from.Location())
	}

	// Final check: ensure from < to
	if to.Before(from) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "'to_datetime' must be after 'from_datetime'"})
	}

	params.Filters["departure_at"] = map[string]interface{}{
		"from": from,
		"to":   to,
	}

	params.Preloads = append(params.Preloads, "Driver")

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

	dtoRides, err := dto.MapToDTOs[models.Ride, dto.RideListResponseDTO](matchingRides)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to map to DTO"})
	}
	// Update response data with matched rides
	result.Data = dtoRides

	return c.JSON(http.StatusOK, result)
}
