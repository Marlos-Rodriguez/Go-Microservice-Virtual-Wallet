package storage

import (
	"github.com/jinzhu/gorm"
)

//AuthStorageService struct
type AuthStorageService struct {
	db *gorm.DB
}

//NewAuthStorageService Return a new Auth Storage Service
func NewAuthStorageService(DB *gorm.DB) *AuthStorageService {
	return &AuthStorageService{db: DB}
}
