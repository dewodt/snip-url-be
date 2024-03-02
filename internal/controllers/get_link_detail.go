package controllers

import (
	"net/http"
	"snip-url-be/internal/auth"
	"snip-url-be/internal/db"
	"snip-url-be/internal/models"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomPath struct {
	ID        uuid.UUID `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
}

type Last4Weeks struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type Referrer struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type Device struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}

type Country struct {
	Country    string  `json:"country"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type LinkDetailResponse struct {
	ID                    uuid.UUID    `json:"id"`
	Title                 string       `json:"title"`
	DestinationUrl        string       `json:"destinationUrl"`
	CreatedAt             time.Time    `json:"createdAt"`
	CustomPaths           []CustomPath `json:"customPaths"`
	Last4Weeks            []Last4Weeks `json:"last4Weeks"`
	Referrers             []Referrer   `json:"referrers"`
	Devices               []Device     `json:"devices"`
	Countries             []Country    `json:"countries"`
	TotalRequests         int64        `json:"totalRequests"`
	LastWeekTotalRequests int64        `json:"lastWeekTotalRequests"`
}

func GetLinkDetailHandler(c *gin.Context) {
	// Get ID from param
	linkID := c.Param("id")

	// Get user from context
	session := auth.GetSessionFromContext(c)

	// Query db
	var link models.Link
	err := db.DB.Preload("CustomPaths").Preload("Requests").Where("user_id = ? AND id = ?", session.ID, linkID).First(&link).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}

	// Get last 4 weeks data
	var last4Weeks []Last4Weeks
	for i := 0; i < 28; i++ {
		date := time.Now().AddDate(0, 0, -i)
		last4Weeks = append(last4Weeks, Last4Weeks{Date: date.Format("2006-01-02"), Count: 0})
		for _, req := range link.Requests {
			if req.CreatedAt.Format("2006-01-02") == date.Format("2006-01-02") {
				last4Weeks[i].Count++
			}
		}
	}

	// Get total count
	total := int64(len(link.Requests))

	// Last week
	lastWeek := last4Weeks[:7]
	totalLastWeeks := int64(0)
	for _, week := range lastWeek {
		totalLastWeeks += week.Count
	}

	// Referrers Map
	referrersMap := make(map[string]int64)
	for _, req := range link.Requests {
		referrersMap[req.Referrer]++
	}
	// Referrers
	var referrers []Referrer
	countNotOtherReferrer := int64(0)
	for k, v := range referrersMap {
		referrers = append(referrers, Referrer{Referrer: k, Count: v})
		countNotOtherReferrer += v
	}
	sort.Slice(referrers, func(i, j int) bool {
		return referrers[i].Count > referrers[j].Count
	})
	referrers = referrers[:min(len(referrers), 5)]
	// Other referrers
	if (total - countNotOtherReferrer) > 0 {
		referrers = append(referrers, Referrer{Referrer: "Other", Count: total - countNotOtherReferrer})
	}

	// Devices Map
	var devicesMap = make(map[string]int64)
	for _, req := range link.Requests {
		devicesMap[req.Device]++
	}
	// Devices
	var devices []Device
	if devicesMap["Desktop"] > 0 {
		devices = append(devices, Device{Device: "Desktop", Count: devicesMap["Desktop"]})
	}
	if devicesMap["Mobile"] > 0 {
		devices = append(devices, Device{Device: "Mobile", Count: devicesMap["Mobile"]})
	}
	if devicesMap["Tablet"] > 0 {
		devices = append(devices, Device{Device: "Tablet", Count: devicesMap["Tablet"]})
	}
	if devicesMap["Unknown"] > 0 {
		devices = append(devices, Device{Device: "Unknown", Count: devicesMap["Unknown"]})
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Count > devices[j].Count
	})

	// Count countries
	countryMap := make(map[string]int64)
	for _, req := range link.Requests {
		countryMap[req.Country]++
	}
	// Countries
	var countries []Country
	countNotOtherCountry := int64(0)
	for k, v := range countryMap {
		percentage := float64(v) / float64(total) * 100
		countries = append(countries, Country{Country: k, Count: v, Percentage: percentage})
		countNotOtherCountry += v
	}
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Count > countries[j].Count
	})
	countries = countries[:min(len(countries), 5)]
	// Other countries
	if (total - countNotOtherCountry) > 0 {
		otherPercentage := float64(total-countNotOtherCountry) / float64(total) * 100
		countries = append(countries, Country{Country: "Other", Count: total - countNotOtherCountry, Percentage: otherPercentage})
	}

	// Custom paths
	var customPaths = make([]CustomPath, len(link.CustomPaths))
	for i, cp := range link.CustomPaths {
		customPaths[i] = CustomPath{
			ID:        cp.ID,
			Path:      cp.Path,
			CreatedAt: cp.CreatedAt,
		}
	}

	// Return link
	c.JSON(http.StatusOK, gin.H{"data": LinkDetailResponse{
		ID:                    link.ID,
		Title:                 link.Title,
		CreatedAt:             link.CreatedAt,
		CustomPaths:           customPaths,
		DestinationUrl:        link.DestinationUrl,
		Last4Weeks:            last4Weeks,
		Referrers:             referrers,
		Devices:               devices,
		Countries:             countries,
		TotalRequests:         total,
		LastWeekTotalRequests: totalLastWeeks,
	}})
}
