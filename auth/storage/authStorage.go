package storage

import (
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	grpcClient "github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/utils"
	UserModels "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//AuthStorageService struct
type AuthStorageService struct {
	db  *gorm.DB
	rdb *redis.Client
}

//NewAuthStorageService Return a new Auth Storage Service
func NewAuthStorageService(DB *gorm.DB, RDB *redis.Client) *AuthStorageService {
	DB.AutoMigrate(&UserModels.User{}, &UserModels.Profile{})

	grpcClient.StartMoveClient()

	return &AuthStorageService{db: DB, rdb: RDB}
}

//CloseDB Close Postgres DB and Redis DB
func (s *AuthStorageService) CloseDB() {
	s.db.Close()
	s.rdb.Close()
	grpcClient.CloseMoveClient()
}

func (s *AuthStorageService) register(newUser *UserModels.User) (bool, error) {
	//Check if if have all requery paramts
	if newUser.UserName == "" || len(newUser.UserName) <= 0 || newUser.Profile.Email == "" || len(newUser.Profile.Email) <= 0 {
		return false, errors.New("Review your Input")
	}

	//Check if username & email exits
	exits, err := s.CheckExistingUser(newUser.UserName, newUser.Profile.Email)

	if err != nil {
		return false, err
	}

	if exits {
		return false, errors.New("Username or Email already in use")
	}

	/* Create in DB */

	//Create User
	newUser.UserID = uuid.New()
	newUser.Balance = 0
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	newUser.IsActive = true

	//Create Profile
	encryptPassword, err := utils.EncryptPassword(newUser.Profile.Password)
	if err != nil {
		return false, nil
	}

	newUser.Profile.UserID = newUser.UserID
	newUser.Profile.Password = encryptPassword
	newUser.Profile.CreatedAt = time.Now()
	newUser.Profile.UpdatedAt = time.Now()
	newUser.Profile.IsActive = true

	go s.db.Create(&newUser)
	s.db.Create(&newUser.Profile)

	if err := s.db.Error; err != nil {
		return false, err
	}

	//Create in Cache
	s.SetRegisterCache(newUser.UserName, newUser.Profile.Email)

	//Create movement
	change := "New User with UserName " + newUser.UserName + "and Email " + newUser.Profile.Email

	success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

	if !success || err != nil {
		log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
	}

	return true, nil
}
