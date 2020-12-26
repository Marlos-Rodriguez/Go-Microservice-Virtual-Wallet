package storage

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
)

//SetTransactionCache Set transactions in Cache
func (s *TransactionStorageService) SetTransactionCache(id string, transactions []*models.TransactionResponse) {
	transactionsCache, _ := json.Marshal(transactions)

	var key string = "Transactions:" + id

	status := s.rdb.Set(context.Background(), key, transactionsCache, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//GetTransactionsCache Get transactions save in Cache
func (s *TransactionStorageService) GetTransactionsCache(id string) ([]*models.TransactionResponse, error) {
	//Get info from redis
	val := s.rdb.Get(context.Background(), "Transactions:"+id)

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

//UpdateTransactionCache Update the last transactions in the cache
func (s *TransactionStorageService) UpdateTransactionCache(id string) {
	//Get in DB
	var transactionsDB []*models.Transaction = []*models.Transaction{new(models.Transaction)}

	limit := 30

	if err := s.db.Order("created_at desc").Where("from_user = ?", id).Or("to_user = ?", id).Find(&transactionsDB).Limit(limit).Error; err != nil {
		log.Println("Error in Update transaction Cache " + err.Error())
	}

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

	s.SetTransactionCache(id, transactionsResponse)
}
