package storage

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	grpcClient "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
)

//TransactionStorageService struct
type TransactionStorageService struct {
	db  *gorm.DB
	rdb *redis.Client
}

//NewTransactionStorageService Create a new storage service
func NewTransactionStorageService(db *gorm.DB, rdb *redis.Client) TransactionStorageService {
	go grpcClient.StartClient()

	db.AutoMigrate(&models.Transaction{})

	newService := TransactionStorageService{db: db, rdb: rdb}

	return newService
}

//CloseDB and grpc Client
func (s *TransactionStorageService) CloseDB() {
	grpcClient.CloseClient()
	s.db.Close()
	s.rdb.Close()
}

//GetTransactions of User
func (s *TransactionStorageService) GetTransactions(userID string, page int) ([]*models.TransactionWebResponse, error) {
	if userID == "" || len(userID) <= 0 {
		return nil, errors.New("Must send ID")
	}

	//IF is the page is 1, check in cache
	if page <= 1 {
		transactionsCache, err := s.GetTransactionsCache(userID)

		if transactionsCache != nil && err == nil {
			return transactionsCache, nil
		}
	}

	//Get in DB
	var transactionsDB []*models.Transaction = []*models.Transaction{new(models.Transaction)}

	limit := page * 30

	if err := s.db.Order("created_at desc").Where("from_user = ?", userID).Or("to_user = ?", userID).Find(&transactionsDB).Limit(limit).Error; err != nil {
		return nil, err
	}

	//response
	var transactionsResponse []*models.TransactionWebResponse

	for _, transaction := range transactionsDB {
		loopTransaction := models.TransactionWebResponse{
			TsID:      transaction.TsID.String(),
			FromUser:  transaction.FromUser.String(),
			FromName:  transaction.FromName,
			ToUser:    transaction.ToUser.String(),
			ToName:    transaction.ToName,
			Amount:    transaction.Amount,
			Message:   transaction.Message,
			CreatedAt: transaction.CreatedAt.String(),
		}

		transactionsResponse = append(transactionsResponse, &loopTransaction)
	}

	s.SetTransactionCache(userID, transactionsResponse)

	return transactionsResponse, nil
}

//CreateTransaction between users
func (s *TransactionStorageService) CreateTransaction(transaction models.TransactionWebRequest) (*models.TransactionWebResponse, bool, error) {

	//Check For User Active & Amount
	if success, err := grpcClient.CheckUserTransaction(
		transaction.FromUser,
		transaction.ToUser,
		transaction.Amount,
		transaction.Password); !success || err != nil {
		return nil, false, err
	}

	//Update Amount
	if success, err := grpcClient.MakeTransaction(transaction.FromUser, transaction.ToUser, transaction.Amount); !success || err != nil {
		return nil, false, err
	}

	fromID, err := uuid.Parse(transaction.FromUser)
	if err != nil {
		return nil, false, errors.New("Error converting the ID in DB")
	}
	toID, err := uuid.Parse(transaction.ToUser)
	if err != nil {
		return nil, false, errors.New("Error converting the ID in DB")
	}

	//Set to DB transaction
	newTransaction := models.Transaction{
		TsID:      uuid.New(),
		FromUser:  fromID,
		FromName:  transaction.FromName,
		ToUser:    toID,
		ToName:    transaction.ToName,
		Amount:    transaction.Amount,
		Message:   transaction.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	//Create relation in DB
	if err := s.db.Create(&newTransaction).Error; err != nil {
		return nil, false, err
	}

	var wg sync.WaitGroup

	//Update Cache
	go func() {
		s.UpdateTransactionCache(transaction.FromUser)
		wg.Done()
	}()
	go func() {
		s.UpdateTransactionCache(transaction.ToUser)
		wg.Done()
	}()
	go func() {
		movement := fmt.Sprintf("Trasaction of %.2f from %s to %s", transaction.Amount, transaction.FromName, transaction.ToName)

		//Create Movement
		if success, err := grpcClient.CreateMovement("User, Profile & Transaction", movement, "Transaction Service"); !success || err != nil {
			log.Println("Error in create Movement in Transaction service")
		}
		wg.Done()
	}()

	TransactionResponse := models.TransactionWebResponse{
		TsID:      newTransaction.TsID.String(),
		FromUser:  newTransaction.FromUser.String(),
		FromName:  newTransaction.FromName,
		ToUser:    newTransaction.ToUser.String(),
		ToName:    newTransaction.ToName,
		Amount:    newTransaction.Amount,
		Message:   newTransaction.Message,
		CreatedAt: newTransaction.CreatedAt.String(),
	}

	wg.Wait()

	return &TransactionResponse, true, nil
}
