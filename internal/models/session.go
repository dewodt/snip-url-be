package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Token     string
	UserID    uint
	User      User
	CreatedAt time.Time
	ExpiresAt time.Time
}
