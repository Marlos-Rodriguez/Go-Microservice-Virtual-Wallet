package models

import (
	"time"
)

//RegisterRequest request struct
type RegisterRequest struct {
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Birthday  time.Time `json:"birthday"`
	Biography string    `json:"biography"`
}
