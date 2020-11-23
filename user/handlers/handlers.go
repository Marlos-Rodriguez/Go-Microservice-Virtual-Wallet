package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/storage"
)

//UserhandlerService struct
type UserhandlerService struct {
	StorageService *storage.UserStorageService
}

//NewUserhandlerService Create new user handler
func NewUserhandlerService(newDB *gorm.DB) *UserhandlerService {
	//return new Handler service
	return &UserhandlerService{
		StorageService: storage.NewUserStorageService(newDB),
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
	if len(body.CurrentUserName) > 0 && len(body.NewUsername) > 0 {
		userDB.UserName = body.CurrentUserName
		newUserName = body.NewUsername
	}

	//Email
	if len(body.Email) > 0 || body.Email != "" {
		userDB.Profile.Email = body.Email
	}

	//Birthday
	userDB.Profile.Birthday = body.Birthday

	//FirstName
	if len(body.FirstName) > 0 || body.FirstName != "" {
		userDB.Profile.FirstName = body.FirstName
	}

	//LastName
	if len(body.LastName) > 0 || body.LastName != "" {
		userDB.Profile.LastName = body.LastName
	}

	//Password
	if len(body.Password) > 0 || body.Password != "" {
		userDB.Profile.Password = body.Password
	}

	//Biography
	if len(body.Biography) > 0 || body.Biography != "" {
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
	if len(body.FromID) < 0 || body.FromID == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from ID"})
	}

	//From Username
	if len(body.FromName) < 0 || body.FromName == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from Username"})
	}

	//From Email
	if len(body.FromEmail) < 0 || body.FromEmail == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from email"})
	}
	//To ID
	if len(body.ToID) < 0 || body.ToID == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to ID"})
	}
	//To Username
	if len(body.ToName) < 0 || body.ToName == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to Username"})
	}
	//To Email
	if len(body.ToEmail) < 0 || body.ToEmail == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to Email"})
	}

	if sucess, err := u.StorageService.AddRelation(body); sucess != true || err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Error in create in DB", "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "message": "Relation created"})
}
