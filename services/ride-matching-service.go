package services

import (
	"carpool-backend/models"
	"carpool-backend/utils"
	"errors"
	"fmt"
)

// MatchRides matches a rider's origin and destination with available rides using polyline decoding
func MatchRides(riderOriginLat, riderOriginLng, riderDestLat, riderDestLng float64, radius float64, availableRides []models.Ride) ([]models.Ride, error) {
	var matchingRides []models.Ride

	for _, ride := range availableRides {
		// Decode the polyline for the ride
		points, err := utils.DecodePolyline(ride.Route)
		if err != nil {
			fmt.Printf("Warning: Failed to decode polyline for ride ID %d: %v\n", ride.ID, err)
			continue // Skip this ride and process the next one
		}

		originMatched, destinationMatched := false, false

		// Optimize: Check fewer points in the polyline to reduce computation
		step := len(points) / 20 // Sample 20 points evenly (avoids checking all)
		if step < 1 {
			step = 1
		}

		// Check proximity of rider's origin and destination to the polyline points
		for i := 0; i < len(points); i += step {
			point := points[i]

			// Check if rider's origin is near any point in the polyline
			if !originMatched && utils.Haversine(point.Lat, point.Lng, riderOriginLat, riderOriginLng) <= radius {
				originMatched = true
			}
			// Check if rider's destination is near any point in the polyline
			if !destinationMatched && utils.Haversine(point.Lat, point.Lng, riderDestLat, riderDestLng) <= radius {
				destinationMatched = true
			}

			// If both origin and destination are matched, add the ride
			if originMatched && destinationMatched {
				matchingRides = append(matchingRides, ride)
				break // No need to check further for this ride
			}
		}
	}

	if len(matchingRides) == 0 {
		return nil, errors.New("no matching rides found")
	}

	return matchingRides, nil
}
