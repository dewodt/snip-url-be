package db

import (
	"fmt"
	"os"
	"snip-url-be/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDB() {
	// Env
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_SSLMODE := os.Getenv("DB_SSLMODE")

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE)
	db, err := gorm.Open(postgres.Open((dsn)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Pooling
	psqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	psqlDB.SetConnMaxIdleTime(10)
	psqlDB.SetMaxOpenConns(100)
	psqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Verification{}, &models.Session{})
	if err != nil {
		panic("failed to migrate schema")
	}

	// Global variable
	DB = db
}
