package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Read an unsigned integer parameter with the given key
func uintParam(c *gin.Context, key string) (uint64, error) {
	value := c.Param(key)
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: expected uint64, got: %s", key, value)
	}
	return i, nil
}

// Get and return bounded query parameters for paging.
// If no query params are found, default values are returned.
func getPageParams(c *gin.Context) (uint64, int) {
	cursor, limit := uint64(0), 10
	if cursorQuery, ok := c.GetQuery("cursor"); ok {
		cursor, _ = strconv.ParseUint(cursorQuery, 10, 64)
	}
	if limitQuery, ok := c.GetQuery("limit"); ok {
		limit, _ = strconv.Atoi(limitQuery)
	}
	return cursor, clamp(limit)
}

// Ensure limit is between 10 and 1000
func clamp(limit int) int {
	if limit >= 10 && limit <= 1000 {
		return limit
	}
	return 10
}
