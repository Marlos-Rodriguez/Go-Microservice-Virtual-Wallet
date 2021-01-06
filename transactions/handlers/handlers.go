package handlers

import (
	"strconv"
	"strings"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/storage"

	internalJWT "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/internal/jwt"
)

//TransactionsHandlerService struct
type TransactionsHandlerService struct {
	storageService storage.TransactionStorageService
}

//NewTransactionsHandlerService Create a new TS handler service
func NewTransactionsHandlerService(newDB *gorm.DB, newRedis *redis.Client) *TransactionsHandlerService {
	return &TransactionsHandlerService{
		storageService: storage.NewTransactionStorageService(newDB, newRedis),
	}
}

//GetTransactions Of User
func (s *TransactionsHandlerService) GetTransactions(c *fiber.Ctx) error {
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

	//Check the JWT ID
	tk := c.Locals("user").(*jwt.Token)
	if err := internalJWT.GetClaims(*tk); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	if match, err := internalJWT.CheckID(ID); !match || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	transactions, err := s.storageService.GetTransactions(ID, pageInt)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in DB", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(transactions)
}

//CreateTransaction between Users
func (s *TransactionsHandlerService) CreateTransaction(c *fiber.Ctx) error {
	var newTransaction *models.TransactionRequest

	if err := c.BodyParser(&newTransaction); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your body", "data": err.Error()})
	}

	//From ID
	if len(strings.TrimSpace(newTransaction.FromUser)) < 0 || strings.TrimSpace(newTransaction.FromUser) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from ID"})
	}

	//Check the JWT ID
	tk := c.Locals("user").(*jwt.Token)
	if err := internalJWT.GetClaims(*tk); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	if match, err := internalJWT.CheckID(newTransaction.FromUser); !match || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
	}

	//From Username
	if len(strings.TrimSpace(newTransaction.FromName)) < 0 || strings.TrimSpace(newTransaction.FromName) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending from Username"})
	}
	newTransaction.FromName = strings.TrimSpace(strings.ToLower(newTransaction.FromName))

	//Password
	if len(strings.TrimSpace(newTransaction.Password)) < 0 || strings.TrimSpace(newTransaction.Password) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending Password"})
	}

	//To ID
	if len(strings.TrimSpace(newTransaction.ToUser)) < 0 || strings.TrimSpace(newTransaction.ToUser) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending to ID"})
	}

	//To Username
	if len(strings.TrimSpace(newTransaction.ToName)) < 0 || strings.TrimSpace(newTransaction.ToName) == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending To Username"})
	}
	newTransaction.ToName = strings.TrimSpace(strings.ToLower(newTransaction.ToName))

	//Amount
	if newTransaction.Amount <= 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error sending Amount"})
	}

	transactionResponse, success, err := s.storageService.CreateTransaction(*newTransaction)

	if !success || err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in Create in DB", "data": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(transactionResponse)
}
