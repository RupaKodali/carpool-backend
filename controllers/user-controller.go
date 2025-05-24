package controllers

import (
	"carpool-backend/dto"
	"carpool-backend/models"
	"carpool-backend/services"
	"carpool-backend/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	UserService services.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// RegisterUser handles POST /users/register
func (h *UserController) RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	count, err := h.UserService.CountUsersByEmailOrPhone(user.Email, user.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check existing user"})
	}
	if count > 0 {
		return c.JSON(http.StatusOK, echo.Map{"error": "Email or phone already in use"})
	}

	count, err = h.UserService.CountUsersByUsername(user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check existing username"})
	}
	if count > 0 {
		return c.JSON(http.StatusOK, echo.Map{"isAvaialble": false, "message": "Username already in use"})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
	}
	user.Password = hashedPassword
	user.AuthProvider = "email"

	if err := h.UserService.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

// LoginUser handles POST /users/login
func (h *UserController) LoginUser(c echo.Context) error {
	var loginRequest struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, err := h.UserService.GetUserByIdentifier(loginRequest.Identifier)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	if err := utils.CheckPassword(user.Password, loginRequest.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	tokens, err := h.GenerateTokens(user.ID, user.IsDriver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	var userResponse dto.UserLoginResponse
	err = copier.Copy(&userResponse, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to copy user to DTO"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":       "Login successful",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user":          userResponse,
	})
}

// RefreshToken handles POST /auth/refresh
func (h *UserController) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil || req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid refresh token"})
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired refresh token"})
	}

	userID := uint(claims["user_id"].(float64))
	accessToken, err := utils.GenerateAccessToken(userID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate new access token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": accessToken})
}

func (h *UserController) GoogleLogin(c echo.Context) error {
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := c.Bind(&req); err != nil || req.IDToken == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing ID token"})
	}

	payload, err := utils.ValidateGoogleIDToken(req.IDToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid ID token"})
	}
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email not found in ID token"})
	}
	name := payload.Claims["name"].(string)

	firstname := payload.Claims["given_name"].(string)
	lastname := payload.Claims["given_name"].(string)

	sub := payload.Claims["sub"].(string) // Google user ID

	// 1. Check if user exists by email
	// 2. If not, create a new user
	user, err := h.UserService.GetUserByEmail(email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error checking user existence"})
		}
		user = &models.User{
			Email:        email,
			FirstName:    firstname,
			LastName:     lastname,
			Username:     name,
			GoogleID:     &sub,
			AuthProvider: "google",
		}
		if err := h.UserService.CreateUser(user); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "User creation failed"})
		}
	}

	// 3. Return JWT access/refresh token from your system
	tokens, err := h.GenerateTokens(user.ID, user.IsDriver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	var userResponse dto.UserLoginResponse
	err = copier.Copy(&userResponse, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to copy user to DTO"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":       "Login successful",
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user":          userResponse,
	})
}

// GetUser handles GET /users/:id
func (h *UserController) GetUser(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	user, err := h.UserService.GetUserByID(int(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id
func (h *UserController) UpdateUser(c echo.Context) error {
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}
	id := uint(id64)

	if loggedInUserID != id {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to update this user"})
	}

	user, err := h.UserService.GetUserByID(int(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
		}
		updates["password"] = hashedPassword
	}

	err = h.UserService.UpdateUser(user, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /users/:id
func (h *UserController) DeleteUser(c echo.Context) error {
	loggedInUserID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}
	id := uint(id64)

	if loggedInUserID != id {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You are not authorized to delete this user"})
	}

	err = h.UserService.DeleteUser(int(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}

func (h *UserController) GenerateTokens(user_id uint, isDriver bool) (tokens dto.TokenStruct, err error) {
	accessToken, err := utils.GenerateAccessToken(user_id, isDriver)
	if err != nil {
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user_id, isDriver)
	if err != nil {
		return
	}

	// TODO: Optional: store refresh token in DB (not implemented here)
	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

// CheckUniqueUsername handles GET /users/:username
func (h *UserController) CheckUniqueUsername(c echo.Context) error {
	username := c.Param("username")

	count, err := h.UserService.CountUsersByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check existing username"})
	}
	if count > 0 {
		return c.JSON(http.StatusOK, echo.Map{"isAvaialble": false, "message": "Username already in use"})
	}

	return c.JSON(http.StatusOK, echo.Map{"isAvaialble": true, "message": "Username available"})
}

// RefreshToken handles POST /auth/forgot-password
func (h *UserController) ForgotPassword(c echo.Context) error {
	var req struct {
		Identifier string `json:"identifier"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	//email and sms otp code

	Otp := "1234"
	updates := map[string]interface{}{"otp": Otp}

	err := h.UserService.UpdateUserByIdentifier(req.Identifier, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate OTP"})
	}

	return c.JSON(http.StatusOK, echo.Map{"otp": Otp, "message": "OTP sent successfully"})
}

// ValidateOtp handles POST /auth/validate-otp
func (h *UserController) ValidateOtp(c echo.Context) error {
	var req struct {
		Identifier string `json:"identifier"`
		Otp        string `json:"otp"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, err := h.UserService.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}
	IsOtpVerified := false

	if req.Otp == user.Otp {
		IsOtpVerified = true
	}
	return c.JSON(http.StatusOK, echo.Map{"is_otp_verified": IsOtpVerified, "message": "OTP verified successfully"})
}

// UpdatePassword handles POST /auth/update-password
func (h *UserController) UpdatePassword(c echo.Context) error {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, err := h.UserService.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	if user.Otp == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to set Password! Please try again!"})

	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
	}

	updates := map[string]interface{}{"password": hashedPassword, "otp": ""}

	err = h.UserService.UpdateUserByIdentifier(req.Identifier, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to set Password! Please try again!"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Password set successfully"})
}
