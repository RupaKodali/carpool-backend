package services

import (
	"carpool-backend/configs"
	"context"
	"log"

	"googlemaps.github.io/maps"
)

func GetCoordinates(address string) (float64, float64, error) {
	apiKey := configs.GetGoogleMapsAPIKey()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating Google Maps client: %v", err)
	}

	req := &maps.GeocodingRequest{
		Address: address,
	}

	results, err := client.Geocode(context.Background(), req)
	if err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, nil // No results found
	}

	location := results[0].Geometry.Location
	return location.Lat, location.Lng, nil
}
