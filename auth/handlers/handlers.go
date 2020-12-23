package handlers

import (
	"strings"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/storage"
	internalJWT "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/jwt"
	UserModels "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

//AuthHandlerService Handler Service struct
type AuthHandlerService struct {
	storageService *storage.AuthStorageService
}

//NewAuthHandlerService Create a new Auth Handler Service
func NewAuthHandlerService(DB *gorm.DB, RDB *redis.Client) AuthHandlerService {
	newService := AuthHandlerService{storageService: storage.NewAuthStorageService(DB, RDB)}
	return newService
}

//Register a new User
func (s *AuthHandlerService) Register(c *fiber.Ctx) error {
	var newUser models.RegisterRequest

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	/* Check all the Values */

	var userDB UserModels.User

	/* Essential */

	//Username
	if len(strings.TrimSpace(newUser.Username)) < 0 || strings.TrimSpace(newUser.Username) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userDB.UserName = strings.ToLower(strings.TrimSpace(newUser.Username))

	//Email
	if len(strings.TrimSpace(newUser.Email)) < 0 || strings.TrimSpace(newUser.Email) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Email"})
	}
	userDB.Profile.Email = strings.ToLower(strings.TrimSpace(newUser.Email))

	//Password
	if len(newUser.Password) < 6 || newUser.Password == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "The password must be 6 character"})
	}
	userDB.Profile.Password = newUser.Password

	/* Optionals */

	//First Name
	if len(strings.TrimSpace(newUser.FirstName)) > 0 || strings.TrimSpace(newUser.FirstName) != "" {
		userDB.Profile.FirstName = newUser.FirstName
	}

	//Last Name
	if len(strings.TrimSpace(newUser.LastName)) > 0 || strings.TrimSpace(newUser.LastName) != "" {
		userDB.Profile.LastName = newUser.LastName
	}

	//Biografy
	if len(strings.TrimSpace(newUser.Biography)) > 0 || strings.TrimSpace(newUser.Biography) != "" {
		userDB.Profile.Biography = newUser.Biography
	}

	//Birthday
	if date, err := time.Parse("2006-01-02", newUser.Birthday); err == nil {
		userDB.Profile.Birthday = date
	}

	if success, err := s.storageService.Register(&userDB); err != nil || !success {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB creating the User", "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "message": "User Created created"})
}

//Login the User
func (s *AuthHandlerService) Login(c *fiber.Ctx) error {
	var userBody models.LoginRequest

	if err := c.BodyParser(&userBody); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//Username
	if len(strings.TrimSpace(userBody.Username)) < 0 || strings.TrimSpace(userBody.Username) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Username = strings.ToLower(strings.TrimSpace(userBody.Username))

	//Email
	if len(strings.TrimSpace(userBody.Email)) < 0 || strings.TrimSpace(userBody.Email) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Email = strings.ToLower(strings.TrimSpace(userBody.Email))

	//Password
	if len(strings.TrimSpace(userBody.Password)) < 0 || strings.TrimSpace(userBody.Password) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}

	//Login in DB
	userClaims, success, err := s.storageService.Login(&userBody)

	if !success || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error login in DB", "data": err.Error()})
	}

	//Create JWT
	newToken, err := genereateJWT(*userClaims)

	if newToken == "" || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in create JWT", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted", "token": newToken})
}

//ReactivateUser Active a User
func (s *AuthHandlerService) ReactivateUser(c *fiber.Ctx) error {
	var userBody models.LoginRequest

	if err := c.BodyParser(&userBody); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//Username
	if len(strings.TrimSpace(userBody.Username)) < 0 || strings.TrimSpace(userBody.Username) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Username = strings.ToLower(strings.TrimSpace(userBody.Username))

	//Email
	if len(strings.TrimSpace(userBody.Email)) < 0 || strings.TrimSpace(userBody.Email) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Email = strings.ToLower(strings.TrimSpace(userBody.Email))

	//Password
	if len(strings.TrimSpace(userBody.Password)) < 0 || strings.TrimSpace(userBody.Password) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}

	//Update in DB
	if sucess, err := s.storageService.ReactivateUser(&userBody); !sucess || err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in Update User", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted", "message": "User Reactivated"})
}

//DeactivateUser Deactivate a User
func (s *AuthHandlerService) DeactivateUser(c *fiber.Ctx) error {
	var userBody models.DeactivateUserRequest

	if err := c.BodyParser(&userBody); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//ID
	if len(strings.TrimSpace(userBody.ID)) < 0 || strings.TrimSpace(userBody.ID) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send ID"})
	}

	//Check the JWT ID
	tk := c.Locals("user").(*jwt.Token)
	if err := internalJWT.GetClaims(*tk); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	if match, err := internalJWT.CheckID(userBody.ID); !match || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	//Username
	if len(strings.TrimSpace(userBody.Username)) < 0 || strings.TrimSpace(userBody.Username) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Username = strings.ToLower(strings.TrimSpace(userBody.Username))

	//Email
	if len(strings.TrimSpace(userBody.Email)) < 0 || strings.TrimSpace(userBody.Email) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}
	userBody.Email = strings.ToLower(strings.TrimSpace(userBody.Email))

	//Password
	if len(strings.TrimSpace(userBody.Password)) < 0 || strings.TrimSpace(userBody.Password) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Must Send Username"})
	}

	//Update in DB
	if sucess, err := s.storageService.DeactivateUser(userBody); !sucess || err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in Update User", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted", "message": "User Deleted"})
}
