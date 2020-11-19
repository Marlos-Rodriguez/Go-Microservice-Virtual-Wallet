package models

import (
	"time"

	"github.com/google/uuid"
)

//Movement struct for every change in DB
type Movement struct {
	MovementID uuid.UUID `gorm:"unique_index;not null;type:uuid;default:uuid_generate_v4();primaryKey"`
	Relation   string    `gorm:"not null"`
	Change     string    `gorm:"not null"`
	Origin     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
}
