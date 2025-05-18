package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GenerateToken generates a JWT token for a given user ID
func GenerateAccessToken(userID uint, isDriver bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration: 24 hours
		"isDriver": isDriver,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString(jwtSecret)
}

// GenerateRefreshToken - long-lived (e.g., 7–30 days)
func GenerateRefreshToken(userID uint, isDriver bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"isDriver": isDriver,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}

// ValidateRefreshToken parses and validates a refresh token
func ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to parse token claims")
	}

	return claims, nil
}

// GetUserIDFromToken extracts the user ID from the JWT token in the request context
func GetUserIDFromToken(c echo.Context) (uint, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("unable to extract user ID from token")
	}
	return uint(userIDFloat), nil
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
