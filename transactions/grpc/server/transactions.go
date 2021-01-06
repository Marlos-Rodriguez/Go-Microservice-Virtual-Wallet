package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/internal/cache"
	internal "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/storage"
)

//Server User Server struct
type Server struct {
}

var storageService storage.TransactionStorageService

//GetStorageService Start the storage service for GPRC server
func GetStorageService() {
	db := internal.ConnectDB()
	rDB := cache.NewRedisClient()

	storageService = storage.NewTransactionStorageService(db, rDB)
}

//CloseDB Close both DB
func CloseDB() {
	storageService.CloseDB()
}

//GetTransactions of User
func (s *Server) GetTransactions(ctx context.Context, request *GetTransactionRequest) (*LastTransactionsResponse, error) {
	var response []*Transaction
	if len(request.ID) < 0 || request.ID == "" {
		return &LastTransactionsResponse{Transactions: response}, errors.New("Must send a ID")
	}

	tsDB, err := storageService.GetTransactions(request.ID, 1)

	if err != nil {
		return &LastTransactionsResponse{Transactions: response}, errors.New("Must send a ID")
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

	return &LastTransactionsResponse{Transactions: response}, nil
}
