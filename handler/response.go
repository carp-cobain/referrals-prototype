package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Sends a 200 JSON response
func okJson(c *gin.Context, body gin.H) {
	c.JSON(http.StatusOK, body)
}

// Sends a 201 JSON response
func createdJson(c *gin.Context, body gin.H) {
	c.JSON(http.StatusCreated, body)
}

// Sends an error JSON response.
func errorJson(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// Sends a 400 error JSON response.
func badRequestJson(c *gin.Context, err error) {
	errorJson(c, http.StatusBadRequest, err)
}

// Sends a 404 error JSON response.
func notFoundJson(c *gin.Context, err error) {
	errorJson(c, http.StatusNotFound, err)
}
