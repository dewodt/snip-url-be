package controllers

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"snip-url-be/auth"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type UploadAvatarSchema struct {
	Image *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadAvatarHandler(c *gin.Context) {
	// Get file blob from formdata
	var formData UploadAvatarSchema
	bindErr := c.ShouldBind(&formData)
	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Validate file type & size (5MB)
	if formData.Image.Size > 5*1024*1024 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "File size too large"})
		return
	}

	// Validate file type
	if formData.Image.Header.Get("Content-Type") != "image/jpeg" && formData.Image.Header.Get("Content-Type") != "image/png" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Upload to cloudinary
	CLOUD_NAME := os.Getenv("CLOUDINARY_CLOUD_NAME")
	API_KEY := os.Getenv("CLOUDINARY_API_KEY")
	API_SECRET := os.Getenv("CLOUDINARY_API_SECRET")
	cld, cldErr := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	if cldErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to cloudinary"})
		return
	}

	file, err := formData.Image.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Convert file to byte
	fileBytes := make([]byte, formData.Image.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Convert to base64
	fileType := formData.Image.Header.Get("Content-Type")
	b64 := base64.StdEncoding.EncodeToString(fileBytes)
	b64Formatted := fmt.Sprintf("data:%s;base64,%s", fileType, b64)

	// Get user id
	session := auth.GetSessionFromContext(c)

	// Upload to cloudinary
	folderName := "snip-url/user/"
	publicId := session.ID
	overwrite := true
	res, err := cld.Upload.Upload(c, b64Formatted, uploader.UploadParams{
		Folder:    folderName,
		PublicID:  publicId,
		Overwrite: &overwrite,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Return secure URL
	c.JSON(http.StatusOK, gin.H{"message": "Sucess upload image", "data": res.SecureURL})
}
