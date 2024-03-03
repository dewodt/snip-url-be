package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    uuid.UUID
}
