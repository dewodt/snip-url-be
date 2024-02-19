package controllers

import (
	"encoding/json"
	"net/http"
	"snip-url-be/internal/auth"
	"snip-url-be/internal/db"
	"snip-url-be/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomPathResponse struct {
	ID        uuid.UUID `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
}

type LinkResponse struct {
	ID             uuid.UUID            `json:"id"`
	Title          string               `json:"title"`
	DestinationUrl string               `json:"destinationUrl"`
	CreatedAt      time.Time            `json:"createdAt"`
	RequestCount   int64                `json:"requestCount"`
	CustomPaths    []CustomPathResponse `json:"customPaths"`
}

func GetAllLinksHandler(c *gin.Context) {
	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Get all links & requests data
	var data []models.Link
	err := db.DB.Preload("CustomPaths").Preload("Requests").Where("user_id = ?", session.ID).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}

	// Parse data to response
	var res []LinkResponse
	for _, link := range data {
		// Parse custom paths
		var customPaths []CustomPathResponse
		for _, customPath := range link.CustomPaths {
			customPaths = append(customPaths, CustomPathResponse{
				ID:        customPath.ID,
				Path:      customPath.Path,
				CreatedAt: customPath.CreatedAt,
			})
		}

		// Append to response
		link := LinkResponse{
			ID:             link.ID,
			Title:          link.Title,
			DestinationUrl: link.DestinationUrl,
			CreatedAt:      link.CreatedAt,
			RequestCount:   int64(len(link.Requests)),
			CustomPaths:    customPaths,
		}
		res = append(res, link)
	}

	// Convert to api response camelCase
	byte, err := json.Marshal(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing data"})
		return
	}
	var apiResponse []LinkResponse
	json.Unmarshal(byte, &apiResponse)

	// Return links
	c.JSON(http.StatusOK, apiResponse)
}
