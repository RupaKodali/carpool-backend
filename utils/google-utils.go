package utils

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

// ValidateGoogleID verifies the token and returns the payload
func ValidateGoogleIDToken(idToken string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}
	return payload, nil
}
