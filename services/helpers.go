package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// QueryParams defines the structure for filtering, sorting, and pagination
type QueryParams struct {
	Filters map[string]interface{} // field:value pairs for filtering
	Sort    []SortField            // fields to sort by
	Page    int                    // page number (1-based)
	Limit   int                    // items per page
	Search  string                 // search term for text fields
}

// SortField represents a field to sort by and its direction
type SortField struct {
	Field     string
	Direction string // "ASC" or "DESC"
}

// PaginatedResponse wraps the response with pagination info
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
}

// ParseQueryParams extracts filters, sorting, and pagination from request
func ParseQueryParams(c echo.Context) QueryParams {
	params := QueryParams{
		Filters: make(map[string]interface{}),
		Page:    1,
		Limit:   10,
	}

	// Pagination
	if page := c.QueryParam("page"); page != "" {
		params.Page, _ = strconv.Atoi(page)
	}
	if limit := c.QueryParam("limit"); limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	}

	// Search
	if search := c.QueryParam("search"); search != "" {
		params.Search = search
	}

	// Sorting
	if sort := c.QueryParam("sort"); sort != "" {
		fields := strings.Split(sort, ",")
		for _, field := range fields {
			direction := "ASC"
			if strings.HasPrefix(field, "-") {
				direction = "DESC"
				field = field[1:]
			}
			params.Sort = append(params.Sort, SortField{Field: field, Direction: direction})
		}
	}

	// Extract all other filters dynamically
	for key, values := range c.QueryParams() {
		if key != "page" && key != "limit" && key != "search" && key != "sort" {
			if len(values) > 0 {
				params.Filters[key] = values[0]
			}
		}
	}

	return params
}

// ListEntities is a dynamic function for listing with filters, search, sorting, and pagination
func ListEntities(db *gorm.DB, model interface{}, params QueryParams, searchableFields []string) (*PaginatedResponse, error) {
	query := db.Model(model)

	// Apply Dynamic Filters
	for field, value := range params.Filters {
		switch v := value.(type) {
		case map[string]interface{}:
			// Range Filters (e.g., Date Ranges)
			if from, ok := v["from"]; ok {
				query = query.Where(fmt.Sprintf("%s >= ?", field), from)
			}
			if to, ok := v["to"]; ok {
				query = query.Where(fmt.Sprintf("%s <= ?", field), to)
			}
		default:
			// Special case for Origin - partial match
			if strings.ToLower(field) == "origin" || strings.ToLower(field) == "destination" {
				if strVal, ok := value.(string); ok {
					query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+strVal+"%")
				}
			} else {
				// Normal Field Filters
				query = query.Where(fmt.Sprintf("%s = ?", field), value)
			}
		}
	}

	// Apply Search (Across Multiple Fields)
	if params.Search != "" && len(searchableFields) > 0 {
		searchTerm := "%" + params.Search + "%"
		searchQuery := ""
		searchArgs := []interface{}{}

		for _, field := range searchableFields {
			if searchQuery != "" {
				searchQuery += " OR "
			}
			searchQuery += fmt.Sprintf("%s LIKE ?", field)
			searchArgs = append(searchArgs, searchTerm)
		}

		query = query.Where(searchQuery, searchArgs...)
	}

	// Apply Sorting (Default: created_at DESC)
	if len(params.Sort) > 0 {
		for _, sort := range params.Sort {
			query = query.Order(sort.Field + " " + sort.Direction)
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Get Total Count Before Applying Pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Apply Pagination
	limit := params.Limit
	if limit < 1 {
		limit = 10 // Default limit
	}
	offset := (params.Page - 1) * limit

	// Fetch Results
	if err := query.Limit(limit).Offset(offset).Find(model).Error; err != nil {
		return nil, fmt.Errorf("failed to list records: %v", err)
	}

	// Calculate Total Pages (Prevent division by zero)
	totalPages := 1
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &PaginatedResponse{
		Data:       model,
		Total:      int(total),
		Page:       params.Page,
		TotalPages: totalPages,
	}, nil
}
