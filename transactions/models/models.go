package models

import (
	"time"

	"github.com/google/uuid"
)

//Transaction DB struct
type Transaction struct {
	TsID      uuid.UUID `gorm:"unique_index;not null;type:uuid;default:uuid_generate_v4();primaryKey"`
	FromUser  uuid.UUID `gorm:"not null"`
	FromName  string    `gorm:"not null"`
	ToUser    uuid.UUID `gorm:"not null"`
	ToName    string    `gorm:"not null"`
	Amount    float32   `gorm:"not null"`
	message   string
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"not null;default:false"`
}

//TransactionResponse struct
type TransactionResponse struct {
	TsID      uuid.UUID `json:"tsID"`
	FromUser  uuid.UUID `json:"fromId"`
	FromName  string    `json:"fromName"`
	ToUser    uuid.UUID `json:"toId"`
	ToName    string    `json:"toName"`
	Amount    float32   `json:"amount"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"createAt"`
}
