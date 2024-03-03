package models

import (
	"time"

	"github.com/google/uuid"
)

type CustomPath struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Path      string    `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	LinkID    uuid.UUID
}
