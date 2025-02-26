package controllers

import (
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService services.UserService
}

// NewUserController creates a new UserController with the given UserService
func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// RegisterUser handles POST /users/register
func (h *UserController) RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	err := h.UserService.RegisterUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

// LoginUser handles POST /users/login
func (h *UserController) LoginUser(c echo.Context) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, token, err := h.UserService.LoginUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// GetUser handles GET /users/:id
func (h *UserController) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserController) UpdateUser(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get user ID from request param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	// Ensure the logged-in user is the same as the user being updated
	if loggedInUserID != id {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to update this user"})
	}

	// Fetch the existing user record
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	// Bind request data into a map
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Handle password update (hash before storing)
	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
		}
		updates["password"] = hashedPassword
	}

	// Pass the existing user and updates map to the service
	err = h.UserService.UpdateUser(user, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /users/:id
func (h *UserController) DeleteUser(c echo.Context) error {
	// Extract logged-in user ID from token
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Get user ID from request param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	// Ensure the logged-in user is the same as the user being updated
	if loggedInUserID != id {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to update this user"})
	}

	err = h.UserService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
