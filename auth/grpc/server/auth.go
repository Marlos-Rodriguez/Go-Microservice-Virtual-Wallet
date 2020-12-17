package server

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/cache"
	internal "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/jinzhu/gorm"
)

//Server User Server struct
type Server struct {
}

var storageService *storage.AuthStorageService

//GetStorageService Start the storage service for GPRC server
func GetStorageService() {
	var DB *gorm.DB = internal.ConnectDB("USER")
	var RDB *redis.Client = cache.NewRedisClient("USER")

	storageService = storage.NewAuthStorageService(DB, RDB)
}

//CloseDB Close both DB
func CloseDB() {
	storageService.CloseDB()
}

//ChangeAuthCache Change in redis the User's Username or email
func (s *Server) ChangeAuthCache(ctx context.Context, request *NewUserInfo) (*AuthResponse, error) {
	if len(request.OldUsername) > 0 && len(request.NewUsername) > 0 {
		success, err := storageService.ChangeRegisterCache(request.OldUsername, request.NewUsername, "", "")
		return &AuthResponse{Success: success}, err
	}

	if len(request.OldEmail) > 0 && len(request.NewEmail) > 0 {
		success, err := storageService.ChangeRegisterCache("", "", request.NewEmail, request.OldEmail)
		return &AuthResponse{Success: success}, err
	}

	return &AuthResponse{Success: false}, errors.New("Invalid Input")
}
