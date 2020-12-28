package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/cache"
	internal "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	"github.com/jinzhu/gorm"
)

//Server User Server struct
type Server struct {
}

var db *gorm.DB
var rDB *redis.Client

//GetStorageService Start the storage service for GPRC server
func GetStorageService() {
	db = internal.ConnectDB("TS")
	rDB = cache.NewRedisClient("TS")
}

//CloseDB Close both DB
func CloseDB() {
	db.Close()
	rDB.Close()
}

//GetTransactions of User
func (s *Server) GetTransactions(ctx context.Context, request *TransactionRequest) (*TransactionResponse, error) {
	var response []*Transaction
	if len(request.ID) < 0 || request.ID == "" {
		return &TransactionResponse{Transactions: response}, errors.New("Must send a ID")
	}

	tsDB, err := GetTransactionsCache(request.ID)

	if tsDB != nil {
		tsDB, err = GetTransactions(request.ID)

		if err != nil {
			return &TransactionResponse{Transactions: response}, err
		}
	}

	for _, ts := range tsDB {
		loopTS := Transaction{
			TsID:     ts.TsID,
			FromID:   ts.FromUser,
			ToID:     ts.ToUser,
			FromName: ts.FromName,
			ToName:   ts.ToName,
			Amount:   ts.Amount,
			Message:  ts.Message,
			CreateAt: ts.CreatedAt,
		}

		response = append(response, &loopTS)
	}

	return &TransactionResponse{Transactions: response}, nil
}

//GetTransactionsCache Get transactions save in Cache
func GetTransactionsCache(id string) ([]*models.TransactionResponse, error) {
	//Get info from redis
	val := rDB.Get(context.Background(), "Transactions:"+id)

	err := val.Err()

	if err != nil && err != redis.Nil {
		log.Println("Error in get the cache " + err.Error())
		return nil, err
	}

	var transactionsCache []*models.TransactionResponse

	if err != redis.Nil {
		transactionsBytes, err := val.Bytes()
		if err != nil {
			return nil, err
		}

		json.Unmarshal(transactionsBytes, &transactionsCache)

		return transactionsCache, nil
	}

	return nil, errors.New("Not found in the cache")
}

//GetTransactions of User
func GetTransactions(userID string) ([]*models.TransactionResponse, error) {
	if userID == "" || len(userID) <= 0 {
		return nil, errors.New("Must send ID")
	}

	//Get in DB
	var transactionsDB []*models.Transaction = []*models.Transaction{new(models.Transaction)}

	limit := 30

	if err := db.Order("created_at desc").Where("from_user = ?", userID).Or("to_user = ?", userID).Find(&transactionsDB).Limit(limit).Error; err != nil {
		return nil, err
	}

	//response
	var transactionsResponse []*models.TransactionResponse

	for _, transaction := range transactionsDB {
		loopTransaction := models.TransactionResponse{
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

	return transactionsResponse, nil
}
