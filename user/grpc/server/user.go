package server

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/cache"
	internal "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/storage"
	"github.com/jinzhu/gorm"
)

//Server User Server struct
type Server struct {
}

var storageService *storage.UserStorageService

//GetStorageService Start the storage service for GPRC server
func GetStorageService() {
	var DB *gorm.DB = internal.ConnectDB("5432")
	var RDB *redis.Client = cache.NewRedisClient()

	storageService = storage.NewUserStorageService(DB, RDB)
}

//CloseDB Close both DB
func CloseDB() {
	storageService.CloseDB()
}

//CheckUser check if the user exits
func (s *Server) CheckUser(ctx context.Context, request *UserRequest) (*UserResponse, error) {
	if len(request.ID) < 0 || request.ID == "" {
		return &UserResponse{Exits: false, Active: false}, errors.New("Must send a ID")
	}

	exits, isActive, err := storageService.CheckExistingUser(request.ID)

	if err != nil {
		return &UserResponse{Exits: false, Active: false}, err
	}

	storageService.CloseDB()

	return &UserResponse{Exits: exits, Active: isActive}, nil
}

//CheckRelation Check if exits a Relation
func (s *Server) CheckRelation(ctx context.Context, request *RelationRequest) (*RelationResponse, error) {
	if len(request.FromUsername) < 0 || request.FromUsername == "" && len(request.ToUsername) < 0 || request.ToUsername == "" {
		return &RelationResponse{Exits: false}, errors.New("Must send ID")
	}

	exits, err := storageService.CheckExistingRelation(request.FromUsername, request.ToUsername, false)

	if err != nil {
		return &RelationResponse{Exits: false}, err
	}

	storageService.CloseDB()

	return &RelationResponse{Exits: exits}, nil
}

//ChangeAvatar Change the avatar in DB
func (s *Server) ChangeAvatar(ctx context.Context, request *AvatarName) (*AvatarResponse, error) {
	if len(request.Name) < 0 || request.Name == "" {
		return &AvatarResponse{Sucess: false}, errors.New("Must send the avatar name")
	}

	var userDB *models.User = new(models.User)

	userDB.Profile.Avatar = request.Name

	if sucess, err := storageService.ModifyUser(userDB, "", ""); sucess == false || err != nil {
		return &AvatarResponse{Sucess: false}, err
	}

	return &AvatarResponse{Sucess: true}, nil
}

//SetCacheUser Set the new User in Cache for Get
func (s *Server) SetCacheUser(ctx context.Context, request *UserRequest) (*SetCacheResponse, error) {
	if len(request.ID) < 0 || request.ID == "" {
		return &SetCacheResponse{Done: false}, errors.New("Must send the User ID")
	}

	storageService.UpdateUserCache(request.ID)

	return &SetCacheResponse{Done: true}, nil
}
