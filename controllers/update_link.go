package controllers

import (
	"errors"
	"net/http"
	"snip-url-be/auth"
	"snip-url-be/db"
	"snip-url-be/models"
	"snip-url-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UpdateLinkSchema struct {
	Title      string `form:"title" binding:"required"`
	CustomPath string `form:"customPath" binding:"required"`
}

func UpdateLinkHandler(c *gin.Context) {
	// Validate & bind form data
	formData := UpdateLinkSchema{}
	err := c.ShouldBind(&formData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Get param id
	linkID := c.Param("id")

	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Get link data
	var link models.Link
	err = db.DB.Where("id = ? AND user_id", linkID, session.ID).Preload("CustomPaths").First(&link).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid link ID, link not found"})
		return
	}

	// Check if user is authorized to update link
	if link.UserID.String() != session.ID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Check if latest custom path is the same
	isTitleChanged := link.Title != formData.Title
	isCustomPathChanged := link.CustomPaths[len(link.CustomPaths)-1].Path != formData.CustomPath

	// Update link titled
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if isTitleChanged {
			// Update title
			err := tx.Model(&models.Link{}).Where("id = ? AND user_id = ?", linkID, session.ID).Update("title", formData.Title).Error
			if err != nil {
				return err
			}
		}

		if isCustomPathChanged {
			// Add custom path
			userUUID, _ := utils.StringToUUID(session.ID)
			link := models.Link{
				ID:     link.ID,
				UserID: userUUID,
			}
			err := tx.Model(&link).Association("CustomPaths").Append(&models.CustomPath{Path: formData.CustomPath})
			if err != nil {
				return err
			}
		}

		return nil
	})
	// Errors
	if err != nil {
		var psqlErr *pgconn.PgError
		// Check if error is unique constraint violation
		if errors.As(err, &psqlErr) && (psqlErr.Code == "23505") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Custom path already exists", "field": "customPath"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update link"})
			return
		}
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}
