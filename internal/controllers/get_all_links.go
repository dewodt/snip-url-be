package controllers

import (
	"net/http"
	"snip-url-be/internal/auth"
	"snip-url-be/internal/db"
	"snip-url-be/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AllLinksResponse struct {
	ID             uuid.UUID    `json:"id"`
	Title          string       `json:"title"`
	DestinationUrl string       `json:"destinationUrl"`
	CreatedAt      time.Time    `json:"createdAt"`
	CustomPaths    []CustomPath `json:"customPaths"`
	TotalRequests  int64        `json:"totalRequests"`
}

func GetAllLinksHandler(c *gin.Context) {
	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Get all links & requests data
	var data []models.Link
	err := db.DB.Preload("CustomPaths").Preload("Requests").Where("user_id = ?", session.ID).Order("created_at DESC").Find(&data).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}

	// Parse data to response
	var res []AllLinksResponse
	for _, link := range data {
		// Parse custom paths
		var customPaths []CustomPath
		for _, customPath := range link.CustomPaths {
			customPaths = append(customPaths, CustomPath{
				ID:        customPath.ID,
				Path:      customPath.Path,
				CreatedAt: customPath.CreatedAt,
			})
		}

		// Append to response
		link := AllLinksResponse{
			ID:             link.ID,
			Title:          link.Title,
			DestinationUrl: link.DestinationUrl,
			CreatedAt:      link.CreatedAt,
			TotalRequests:  int64(len(link.Requests)),
			CustomPaths:    customPaths,
		}
		res = append(res, link)
	}

	// Return links
	c.JSON(http.StatusOK, res)
}
