package db

import (
	"snip-url-be/internal/models"
	"snip-url-be/internal/utils"
	"time"
)

func SeedDB() {
	// Drop previous tables
	// DB.Migrator().DropTable(&models.Request{})
	// DB.Migrator().DropTable(&models.CustomPath{})
	// DB.Migrator().DropTable(&models.Link{})
	// DB.Migrator().DropTable(&models.Session{})
	// DB.Migrator().DropTable(&models.Verification{})
	// DB.Migrator().DropTable(&models.Session{})
	// DB.Migrator().DropTable(&models.User{})

	// Migrate the schema
	// err := DB.AutoMigrate(&models.User{}, &models.Verification{}, &models.Session{}, &models.Link{}, &models.CustomPath{}, &models.Request{})
	// if err != nil {
	// 	panic("failed to migrate schema")
	// }

	// User
	uuidUser, err := utils.StringToUUID("3a2068df-82de-4a96-a520-fbd409d96bd1")
	if err != nil {
		panic(err)
	}
	avatar := "https://res.cloudinary.com/dvzs47hay/image/upload/v1708167040/snip-url/user/3a2068df-82de-4a96-a520-fbd409d96bd1.jpg"
	users := &models.User{
		ID:     uuidUser,
		Email:  "dewantorotriatmojo@gmail.com",
		Name:   "Dewanto Triatmojo",
		Avatar: &avatar,
	}

	// Link
	uuidLink1, err := utils.StringToUUID("548e484d-7a53-496d-8c54-e66b76549244")
	if err != nil {
		panic(err)
	}
	uuidLink2, err := utils.StringToUUID("17edc4a1-05af-4f66-ae10-b8d928955c0c")
	if err != nil {
		panic(err)
	}

	links := []models.Link{
		{
			ID:             uuidLink1,
			UserID:         uuidUser,
			Title:          "Link 1",
			DestinationUrl: "https://www.google.com",
			CustomPaths: []models.CustomPath{
				{
					Path: "customPath11",
				},
				{
					Path: "customPath12",
				},
				{
					Path: "customPath13",
				},
			},
			Requests: []models.Request{
				{
					Country:  "US",
					Referrer: "https://www.google.com",
					Device:   "Desktop",
				},
				{
					Country:  "US",
					Referrer: "https://www.google.com",
					Device:   "Desktop",
				},
				{
					Country:  "ID",
					Referrer: "https://www.instagram.com",
					Device:   "Mobile",
				},
				{
					Country:  "ENG",
					Referrer: "https://www.instagram.com",
					Device:   "Tablet",
				},
				{
					Country:  "ENG",
					Referrer: "https://www.instagram.com",
					Device:   "Tablet",
				},
				{
					Country:  "ENG",
					Referrer: "https://www.twitter.com",
					Device:   "Tablet",
				},
				{
					Country:  "ID",
					Referrer: "https://www.instagram.com",
					Device:   "Mobile",
				},
			},
		},
		{
			ID:             uuidLink2,
			Title:          "Link 2",
			UserID:         uuidUser,
			DestinationUrl: "https://www.instagram.com",
			CustomPaths: []models.CustomPath{
				{
					Path: "customPath21",
				},
				{
					Path: "customPath22",
				},
				{
					Path: "customPath23",
				},
			},
			Requests: []models.Request{
				{
					Country:   "US",
					Referrer:  "https://www.google.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -1),
				},
				{
					Country:   "US",
					Referrer:  "https://www.google.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -2),
				},
				{
					Country:   "ID",
					Referrer:  "https://www.instagram.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -3),
				},
				{
					Country:   "US",
					Referrer:  "https://www.youtube.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -1),
				},
				{
					Country:   "US",
					Referrer:  "https://www.google.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -28),
				},
				{
					Country:   "US",
					Referrer:  "https://www.google.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -29),
				},
				{
					Country:   "ENG",
					Referrer:  "https://www.twitter.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -30),
				},
				{
					Country:   "US",
					Referrer:  "https://www.google.com",
					Device:    "Desktop",
					CreatedAt: time.Now().AddDate(0, 0, -15),
				},
			},
		},
	}

	DB.Create(users)
	DB.Create(links)
}
