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
		points, err := utils.DecodePolyline(ride.Route)
		if err != nil {
			fmt.Printf("Warning: Failed to decode polyline for ride ID %d: %v\n", ride.ID, err)
			continue
		}

		originIndex := -1
		destinationIndex := -1

		// Optimization: sample up to 20 points

		step := 1
		if len(points) <= 50 {
			step = 1 // full scan for short routes
		} else {
			step = len(points) / 20 // sample for long routes
		}

		for i := 0; i < len(points); i += step {
			point := points[i]

			if originIndex == -1 && utils.Haversine(point.Lat, point.Lng, riderOriginLat, riderOriginLng) <= radius {
				originIndex = i
			}

			if destinationIndex == -1 && utils.Haversine(point.Lat, point.Lng, riderDestLat, riderDestLng) <= radius {
				destinationIndex = i
			}

			// Exit early if both matched and ordered
			if originIndex != -1 && destinationIndex != -1 {
				if originIndex < destinationIndex {
					matchingRides = append(matchingRides, ride)
				}
				break
			}
		}
	}

	if len(matchingRides) == 0 {
		return nil, errors.New("no matching rides found")
	}

	return matchingRides, nil
}
