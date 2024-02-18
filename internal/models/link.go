package models

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title          string
	DestinationUrl string `gorm:"index:destination_url_user_id,unique"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserID         uuid.UUID `gorm:"index:destination_url_user_id,unique"`
	CustomPaths    []CustomPath
	Requests       []Request
}
