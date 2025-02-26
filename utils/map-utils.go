package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
)

const googleMapsAPI = "https://maps.googleapis.com/maps/api/directions/json"

// Route represents a route returned by the Google Maps API
type Route struct {
	Legs []struct {
		StartLocation struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"start_location"`
		EndLocation struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"end_location"`
		Steps []struct {
			StartLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"start_location"`
			EndLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"end_location"`
		} `json:"steps"`
	} `json:"legs"`
}

// GetRoute retrieves a route between two locations using the Google Maps Directions API
func GetRoute(origin, destination string) (*Route, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	encodedOrigin := url.QueryEscape(origin)
	encodedDestination := url.QueryEscape(destination)

	url := fmt.Sprintf("%s?origin=%s&destination=%s&key=%s", googleMapsAPI, encodedOrigin, encodedDestination, apiKey)

	log.Println("Google Maps API URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Maps API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google Maps API returned status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var routeData struct {
		Routes []Route `json:"routes"`
	}
	if err := json.Unmarshal(body, &routeData); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	if len(routeData.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	return &routeData.Routes[0], nil
}

type Point struct {
	Lat float64
	Lng float64
}

func DecodePolyline(encoded string) ([]Point, error) {
	var points []Point
	index := 0
	length := len(encoded)
	lat := 0.0
	lng := 0.0

	for index < length {
		// Decode latitude
		shift := 0
		result := 0
		for {
			if index >= length {
				return nil, fmt.Errorf("invalid encoding: unexpected end of string at index %d", index)
			}
			b := int(encoded[index]) - 63
			if b < -32 || b > 95 {
				return nil, fmt.Errorf("invalid encoding: byte value out of range at index %d", index)
			}
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}

		// Handle latitude value
		if result&1 != 0 {
			lat += float64(^(result >> 1))
		} else {
			lat += float64(result >> 1)
		}

		// Decode longitude
		shift = 0
		result = 0
		for {
			if index >= length {
				return nil, fmt.Errorf("invalid encoding: unexpected end of string at index %d", index)
			}
			b := int(encoded[index]) - 63
			if b < -32 || b > 95 {
				return nil, fmt.Errorf("invalid encoding: byte value out of range at index %d", index)
			}
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}

		// Handle longitude value
		if result&1 != 0 {
			lng += float64(^(result >> 1))
		} else {
			lng += float64(result >> 1)
		}

		points = append(points, Point{
			Lat: lat / 1e5,
			Lng: lng / 1e5,
		})
	}

	return points, nil
}

// Haversine calculates the distance in kilometers between two geographic coordinates
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	// Convert latitude and longitude to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Differences in coordinates
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin(dLon/2), 2)

	c := 2 * math.Asin(math.Sqrt(a))

	// Calculate distance
	distance := earthRadius * c
	// fmt.Println(distance)
	return distance
}
