package controllers

import (
	"errors"
	"net/http"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/models"
	"snip-url-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateLinkRequest struct {
	Title          string `form:"title" binding:"required"`
	DestinationUrl string `form:"destinationUrl" binding:"required,url"`
	CustomPath     string `form:"customPath" binding:"required,excludesall=~0x2C<> 	;:'\"/[]^{}()=+!*@&$?%#0x7C"`
}

type CreateLinkResponseData struct {
	ID uuid.UUID `json:"id"`
}

func CreateLinkHandler(c *gin.Context) {
	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Get request body
	var formData CreateLinkRequest
	err := c.ShouldBind(&formData)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create uuid from session id
	userId, err := utils.StringToUUID(session.ID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to convert user id"})
		return
	}

	// Clean destination url trailing slash
	destinationUrl := utils.CleanURL(formData.DestinationUrl)

	// Custom path
	customPath := models.CustomPath{
		Path: formData.CustomPath,
	}
	// Link
	link := models.Link{
		ID:             uuid.New(),
		Title:          formData.Title,
		DestinationUrl: destinationUrl,
		UserID:         userId,
		CustomPaths:    []models.CustomPath{customPath},
	}

	// Create link
	err = db.DB.Create(&link).Error

	if err != nil {
		var psqlErr *pgconn.PgError
		// Check if error is unique constraint violation
		if errors.As(err, &psqlErr) && (psqlErr.Code == "23505") {
			if psqlErr.ConstraintName == "destination_url_user_id" {
				// Destination url already exists (in user's links)
				c.AbortWithStatusJSON(400, gin.H{"error": "You already have a link with this destination url", "field": "destinationUrl"})
				return
			} else {
				// Custom path already exists
				c.AbortWithStatusJSON(400, gin.H{"error": "Custom path already exists", "field": "customPath"})
				return
			}
		} else {
			// Other errors
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to create link"})
			return
		}
	}

	// Data send back to client
	responseData := CreateLinkResponseData{
		ID: link.ID,
	}

	// Return link
	c.JSON(http.StatusCreated, gin.H{"message": "Success creating link", "data": responseData})
}
