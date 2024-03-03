package models

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email     string
	Name      *string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"default:now() + interval '10 minutes'"`
}
