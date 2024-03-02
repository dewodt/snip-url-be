package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"snip-url-be/internal/db"
	"snip-url-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

type CountryResponse struct {
	Status      string `json:"status"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

func RedirectHandler(c *gin.Context) {
	// Get path from request
	path := c.Param("customPath")

	// Get link by path
	var link models.Link
	dbRes := db.DB.Raw("SELECT * FROM links WHERE id = (SELECT link_id FROM custom_paths WHERE path = ?)", path).Scan(&link)

	// Not found
	if dbRes.RowsAffected == 0 {
		notFoundURL := os.Getenv("FE_URL") + "/not-found"
		c.Redirect(http.StatusTemporaryRedirect, notFoundURL)
		return
	}
	// Other errors
	if dbRes.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dbRes.Error.Error()})
		return
	}

	// Get device type
	var deviceType string
	userAgent := c.Request.UserAgent()
	ua := useragent.Parse(userAgent)
	if ua.Mobile {
		deviceType = "Mobile"
	} else if ua.Tablet {
		deviceType = "Tablet"
	} else if ua.Desktop {
		deviceType = "Desktop"
	} else {
		deviceType = "Unknown"
	}

	// Get country name
	var country string
	requestIP := c.ClientIP()
	endpoint := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode", requestIP)
	res, err := http.Get(endpoint)
	if err != nil {
		country = "Unknown"
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		country = "Unknown"
	}
	var resJSON CountryResponse
	err = json.Unmarshal(body, &resJSON)
	if err != nil {
		country = "Unknown"
	}
	if resJSON.Status == "fail" {
		country = "Unknown"
	} else {
		country = resJSON.Country
	}

	// Get domain referrer
	var referer string
	reqReferer := c.Request.Referer()
	if reqReferer != "" {
		referer = reqReferer
	} else {
		referer = "Direct"
	}

	// Add stats to db
	newRequest := models.Request{
		Country:  country,
		Device:   deviceType,
		Referrer: referer,
		LinkID:   link.ID,
	}
	db.DB.Create(&newRequest)

	// Redirect to link
	c.Redirect(http.StatusPermanentRedirect, link.DestinationUrl)
}
