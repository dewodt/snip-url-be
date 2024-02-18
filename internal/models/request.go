package models

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Country   string
	Referrer  string
	Device    string
	CreatedAt time.Time
	UpdatedAt time.Time
	LinkID    uuid.UUID
}
