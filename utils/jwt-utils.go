package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GenerateToken generates a JWT token for a given user ID
func GenerateToken(userID int, isDriver bool) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration: 24 hours
		"isDriver": isDriver,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GetUserIDFromToken extracts the user ID from the JWT token in the request context
func GetUserIDFromToken(c echo.Context) (int, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("unable to extract user ID from token")
	}
	return int(userID), nil
}

// IsDriverFromToken extracts the isDriver from the JWT token in the request context
func IsDriverFromToken(c echo.Context) (bool, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	// userID, ok := claims["user_id"].(float64)
	isDriver, ok := claims["isDriver"].(bool)
	if !ok {
		return false, errors.New("unable to extract user ID from token")
	}
	return bool(isDriver), nil
}
