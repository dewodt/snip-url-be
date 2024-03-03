package controllers

import (
	"net/http"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/models"
	"strconv"
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

	// Page size limit
	pageSize := 6

	// Page number
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	// Page Offset
	offset := (page - 1) * pageSize

	// Start time
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		start = 0
	}
	startTime := time.UnixMilli(int64(start))
	startDate := startTime.Format("2006-01-02 15:04:05")
	// End time
	end, err := strconv.Atoi(c.Query("end"))
	if err != nil {
		end = int(time.Now().AddDate(0, 0, 1).UnixMilli()) // Add 1 day to end date to include the end date
	}
	endTime := time.UnixMilli(int64(end)).AddDate(0, 0, 1) // Add 1 day to end date to include the end date
	endDate := endTime.Format("2006-01-02 15:04:05")

	// Get total links
	var totalLinks int64
	err = db.DB.Model(&models.Link{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", session.ID, startDate, endDate).Count(&totalLinks).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}

	// Get all links & requests data
	var data []models.Link
	err = db.DB.
		Preload("CustomPaths").
		Preload("Requests").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", session.ID, startDate, endDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&data).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
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
	c.JSON(http.StatusOK, gin.H{"data": res, "totalLinks": totalLinks})
}
