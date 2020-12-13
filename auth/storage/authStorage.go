package storage

import (
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	grpcClient "github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
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
	go grpcClient.StartClient()

	DB.AutoMigrate(&UserModels.User{}, &UserModels.Profile{})

	return &AuthStorageService{db: DB, rdb: RDB}
}

//CloseDB Close Postgres DB and Redis DB
func (s *AuthStorageService) CloseDB() {
	s.db.Close()
	s.rdb.Close()
	grpcClient.CloseClient()
}

//Register Create a new User
func (s *AuthStorageService) Register(newUser *UserModels.User) (bool, error) {
	//Check if username & email exits
	_, exits, err := s.CheckExistingUser(newUser.UserName, newUser.Profile.Email)

	if err != nil {
		return false, err
	}

	if err == errors.New("User is not active") {
		return false, errors.New("User is not active")
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
	s.SetRegisterCache(newUser.UserName, newUser.Profile.Email, newUser)

	//Create movement
	change := "New User with UserName " + newUser.UserName + "and Email " + newUser.Profile.Email

	success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

	if !success || err != nil {
		log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
	}

	return true, nil
}

//Login Login Funtion
func (s *AuthStorageService) Login(user *models.LoginRequest) (*models.JWTLogin, bool, error) {
	//Get ID from cache
	ID, exits, err := s.CheckExistingUser(user.Username, user.Email)

	if err != nil {
		return nil, false, err
	}

	if !exits && len(ID) <= 0 {
		return nil, false, errors.New("User not exits")
	}

	//Get info from cache
	profileCache, err := s.GetProfileCache(ID)

	if err == nil || profileCache != nil {
		if profileCache.IsActive == false {
			return nil, false, errors.New("User is not active")
		}
		//Obtener ambas contraseñas
		passwordBytes := []byte(user.Password)
		passwordDB := []byte(profileCache.Password)

		//Compare passwords
		if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
			return nil, false, err
		}

		loginResponse := &models.JWTLogin{
			ID:       ID,
			Username: user.Username,
			Email:    profileCache.Email,
			Password: profileCache.Password,
		}

		return loginResponse, true, nil
	}

	//Get info from DB
	var profileDB *UserModels.Profile = new(UserModels.Profile)

	if err := s.db.Where("user_id = ? AND email = ?", ID, user.Email).First(&profileDB).Error; err != nil {
		return nil, false, err
	}

	if !profileDB.IsActive {
		return nil, false, errors.New("User not active")
	}

	//Obtener ambas contraseñas
	passwordBytes := []byte(user.Password)
	passwordDB := []byte(profileDB.Password)

	//Compare passwords
	if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
		return nil, false, err
	}

	loginResponse := &models.JWTLogin{
		ID:       ID,
		Username: user.Username,
		Email:    profileDB.Email,
		Password: profileDB.Password,
	}

	return loginResponse, true, nil
}

//ReactivateUser Reactivate the User
func (s *AuthStorageService) ReactivateUser(ID string, password string) (bool, error) {
	//Get info from cache
	profileCache, err := s.GetProfileCache(ID)

	if err != nil && profileCache != nil {
		if profileCache.IsActive {
			return false, errors.New("User is active")
		}
		//Obtener ambas contraseñas
		passwordBytes := []byte(password)
		passwordDB := []byte(profileCache.Password)

		//Compare passwords
		if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
			return false, err
		}

		go s.db.Model(&UserModels.User{}).Where(&UserModels.User{UserID: profileCache.UserID}).Update("is_active", true)
		s.db.Model(&UserModels.Profile{}).Where(&UserModels.Profile{UserID: profileCache.UserID}).Update("is_active", true)

		if s.db.Error != nil {
			return false, s.db.Error
		}

		//DELETE in Cache
		s.DeleteUserCache(profileCache.UserID.String())

		//Create movement
		change := "User with ID " + ID + "Reactivate his account"

		success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

		if !success || err != nil {
			log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
		}

		return true, nil
	}

	//Get info from DB
	var profileDB *UserModels.Profile = new(UserModels.Profile)

	if err := s.db.Where("user_id = ?", ID).First(&profileDB).Error; err != nil {
		return false, err
	}

	//Obtener ambas contraseñas
	passwordBytes := []byte(password)
	passwordDB := []byte(profileDB.Password)

	//Compare passwords
	if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
		return false, err
	}

	//Update in DB
	go s.db.Model(&UserModels.User{}).Where(&UserModels.User{UserID: profileDB.UserID}).Update("is_active", true)
	s.db.Model(&UserModels.Profile{}).Where(&UserModels.Profile{UserID: profileDB.UserID}).Update("is_active", true)

	if s.db.Error != nil {
		return false, s.db.Error
	}

	//Create movement
	change := "User with ID " + ID + "Reactive his account"

	success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

	if !success || err != nil {
		log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
	}

	return true, nil
}

//DeactivateUser Deactive the User
func (s *AuthStorageService) DeactivateUser(ID string, password string) (bool, error) {
	//Get info from cache
	profileCache, err := s.GetProfileCache(ID)

	if err == nil || profileCache != nil {
		//Obtener ambas contraseñas
		passwordBytes := []byte(password)
		passwordDB := []byte(profileCache.Password)

		//Compare passwords
		if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
			return false, err
		}

		go s.db.Model(&UserModels.User{}).Where(&UserModels.User{UserID: profileCache.UserID}).Update("is_active", false)
		s.db.Model(&UserModels.Profile{}).Where(&UserModels.Profile{UserID: profileCache.UserID}).Update("is_active", false)

		if s.db.Error != nil {
			return false, s.db.Error
		}

		//DELETE in Cache
		s.DeleteUserCache(profileCache.UserID.String())

		//Create movement
		change := "User with ID " + ID + "Deactive his account"

		success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

		if !success || err != nil {
			log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
		}

		return true, nil
	}

	var profileDB *UserModels.Profile = new(UserModels.Profile)

	if err := s.db.Where(&UserModels.User{UserID: profileCache.UserID}).First(&profileDB).Error; err != nil {
		return false, err
	}

	//Obtener ambas contraseñas
	passwordBytes := []byte(password)
	passwordDB := []byte(profileDB.Password)

	//Compare passwords
	if err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes); err != nil {
		return false, err
	}

	//Update in DB
	go s.db.Model(&UserModels.User{}).Where(&UserModels.User{UserID: profileDB.UserID}).Update("is_active", false)
	s.db.Model(&UserModels.Profile{}).Where(&UserModels.Profile{UserID: profileDB.UserID}).Update("is_active", false)

	if s.db.Error != nil {
		return false, s.db.Error
	}

	//DELETE in Cache
	s.DeleteUserCache(profileDB.UserID.String())

	//Create movement
	change := "User with ID " + ID + "Deactive his account"

	success, err := grpcClient.CreateMovement("User & Profile", change, "Auth Service")

	if !success || err != nil {
		log.Println("Error in cretate movement in Auth Service Storage, Register Func. Error: " + err.Error())
	}

	return true, nil
}
