package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name              string
	Email             string `gorm:"uniqueIndex"`
	PasswordHash      string
	EncryptionEnabled bool `gorm:"default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
