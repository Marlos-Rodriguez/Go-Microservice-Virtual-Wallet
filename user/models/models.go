package models

import (
	"time"

	TSModels "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	"github.com/google/uuid"
)

//DB models use to create the tables

//User struct
type User struct {
	UserID    uuid.UUID `gorm:"unique_index;not null;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserName  string    `gorm:"unique_index;not null"`
	Balance   float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"not null;default:false"`
	Profile   Profile   `gorm:"not null;foreignkey:UserID;references:UserID"`
}

//Profile struct
type Profile struct {
	UserID    uuid.UUID `gorm:"unique_index;not null"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique_index;not null"`
	Password  string `gorm:"not null"`
	Birthday  time.Time
	Avatar    string
	Biography string
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"not null;default:false"`
}

//Relation struct
type Relation struct {
	RelationID uuid.UUID `gorm:"unique_index;not null;type:uuid;default:uuid_generate_v4();primaryKey"`
	FromUser   uuid.UUID `gorm:"not null"`
	FromName   string    `gorm:"not null"`
	ToUser     uuid.UUID `gorm:"not null"`
	ToName     string    `gorm:"not null"`
	Mutual     bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	IsActive   bool      `gorm:"not null;default:false"`
}

//Web models create to use in the web responses

//UserRequest struct
type UserRequest struct {
	CurrentUserName string    `json:"current_user_name"`
	NewUsername     string    `json:"new_user_name"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	Birthday        time.Time `json:"birthday"`
	Biography       string    `json:"biography"`
}

//RelationRequest struct
type RelationRequest struct {
	FromID    string `json:"from_ID"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToID      string `json:"to_ID"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
}

//UserResponse struct
type UserResponse struct {
	UserID           uuid.UUID                      `json:"userId"`
	UserName         string                         `json:"username"`
	Balance          float64                        `json:"balance"`
	Avatar           string                         `json:"avatar"`
	LastTransactions []TSModels.TransactionResponse `json:"transactions,omitempty"`
}

//UserProfileResponse struct
type UserProfileResponse struct {
	UserID    uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Birthday  time.Time `json:"birthday"`
	Biography string    `json:"biography"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
}

//RelationReponse struct
type RelationReponse struct {
	RelationID uuid.UUID `json:"relationId"`
	FromUser   uuid.UUID `json:"fromId"`
	FromName   string    `json:"fromName"`
	FromEmail  string    `json:"fromEmail"`
	ToUser     uuid.UUID `json:"toUser"`
	ToName     string    `json:"toName"`
	ToEmail    string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
