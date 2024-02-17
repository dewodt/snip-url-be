package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"unique"`
	Name      string
	Avatar    *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
