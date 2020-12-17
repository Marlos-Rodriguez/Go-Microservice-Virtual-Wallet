package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/storage"
)

//UserhandlerService struct
type UserhandlerService struct {
	StorageService storage.UserStorageService
}

//NewUserhandlerService Create new user handler
func NewUserhandlerService(newDB *gorm.DB, newRDB *redis.Client) *UserhandlerService {
	//return new Handler service
	return &UserhandlerService{
		StorageService: storage.NewUserStorageService(newDB, newRDB),
	}
}

//GetUser Get the basic user Info for main page
func (u *UserhandlerService) GetUser(c *fiber.Ctx) error {
	//Get the ID
	ID := c.Params("id")

	if len(ID) < 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	//Here must be check if the ID of the token mach

	//Get the info from DB
	UserInfo, err := u.StorageService.GetUser(ID)

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB", "data": err.Error()})
	}

	//return the info
	return c.Status(fiber.StatusAccepted).JSON(UserInfo)
}

//GetProfileUser Get the profile info for user info page
func (u *UserhandlerService) GetProfileUser(c *fiber.Ctx) error {
	//Get the ID
	ID := c.Params("id")

	if len(ID) < 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	//Here must be check if the ID of the token mach

	//Get the info from DB
	ProfileInfo, err := u.StorageService.GetProfileUser(ID)

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB", "data": err.Error()})
	}

	//return the info
	return c.Status(fiber.StatusAccepted).JSON(ProfileInfo)
}

//ModifyUser modify the User Info
func (u *UserhandlerService) ModifyUser(c *fiber.Ctx) error {
	//Get the ID
	ID := c.Params("id")

	if len(ID) < 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}
	//Decode the body
	var body models.UserRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//Here must check the id if mach with token

	var userDB models.User
	var newUserName string

	//Username
	if len(strings.TrimSpace(body.CurrentUserName)) > 0 && len(strings.TrimSpace(body.NewUsername)) > 0 {
		userDB.UserName = strings.ToLower(strings.TrimSpace(body.CurrentUserName))
		newUserName = strings.ToLower(strings.TrimSpace(body.NewUsername))
	}

	//Email
	if len(strings.TrimSpace(body.Email)) > 0 || body.Email != "" {
		userDB.Profile.Email = strings.ToLower(strings.TrimSpace(body.Email))
	}

	//Birthday
	if date, err := time.Parse("2006-01-02", body.Birthday); err != nil {
		userDB.Profile.Birthday = date
	}

	//FirstName
	if len(strings.TrimSpace(body.FirstName)) > 0 || strings.TrimSpace(body.FirstName) != "" {
		userDB.Profile.FirstName = strings.ToLower(strings.TrimSpace(body.FirstName))
	}

	//LastName
	if len(strings.TrimSpace(body.LastName)) > 0 || strings.TrimSpace(body.LastName) != "" {
		userDB.Profile.LastName = strings.ToLower(strings.TrimSpace(body.LastName))
	}

	//Password
	if len(body.Password) >= 6 || body.Password != "" {
		userDB.Profile.Password = body.Password
	}

	//Biography
	if len(strings.TrimSpace(body.Biography)) > 0 || strings.TrimSpace(body.Biography) != "" {
		userDB.Profile.Biography = body.Biography
	}

	if sucess, err := u.StorageService.ModifyUser(&userDB, ID, newUserName); err != nil || sucess != true {
		return c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB", "data": err.Error()})
	}

	return c.SendStatus(fiber.StatusAccepted)
}

//GetRelations Get relations from DB
func (u *UserhandlerService) GetRelations(c *fiber.Ctx) error {
	//Get the ID
	ID := c.Params("id")

	if len(ID) < 0 || ID == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	//Get the page
	page := c.Params("page")

	if len(page) < 0 || page == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	//Convert to int
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error converting in Integer", "data": err.Error()})
	}

	//Here must check if the id mach with the token

	//Get info from DB
	relations, err := u.StorageService.GetRelations(ID, pageInt)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(relations)
}

//CreateRelation Create a new relation between users
func (u *UserhandlerService) CreateRelation(c *fiber.Ctx) error {
	//Get the relation info
	var body *models.RelationRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//From ID
	if len(strings.TrimSpace(body.FromID)) < 0 || strings.TrimSpace(body.FromID) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from ID"})
	}

	body.FromEmail = strings.ToLower(body.FromID)
	//From Username
	if len(strings.TrimSpace(body.FromName)) < 0 || strings.TrimSpace(body.FromName) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from Username"})
	}

	body.FromName = strings.ToLower(body.FromName)

	//From Email
	if len(body.FromEmail) < 0 || body.FromEmail == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from email"})
	}

	body.FromEmail = strings.ToLower(body.FromEmail)

	//To Username
	if len(strings.TrimSpace(body.ToName)) < 0 || strings.TrimSpace(body.ToName) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to Username"})
	}

	body.ToName = strings.ToLower(body.ToName)
	//To Email
	if len(strings.TrimSpace(body.ToEmail)) < 0 || strings.TrimSpace(body.ToEmail) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to Email"})
	}

	body.ToEmail = strings.ToLower(body.ToEmail)

	if sucess, err := u.StorageService.AddRelation(body); sucess != true || err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Error in create in DB", "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "message": "Relation created"})
}
