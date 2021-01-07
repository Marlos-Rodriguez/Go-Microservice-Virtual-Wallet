package storage

import (
	"context"
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"google.golang.org/grpc"
)

var tsClient models.TransactionServiceClient
var tsConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startTransactionClient() {
	urlTarget := environment.AccessENV("TRANSACTION_GRPC")

	if urlTarget == "" {
		log.Fatalln("Error in Access to TRANSACTION_GRPC URL in User Service")
	}
	moveConn, err := grpc.Dial(urlTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	tsClient = models.NewTransactionServiceClient(moveConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeTransactionClient() {
	tsConn.Close()
}

//GetTransactions of user
func GetTransactions(id string) ([]models.TransactionResponse, bool, error) {
	transactionsRequest := &models.GetTransactionRequest{
		ID: id,
	}

	response, err := tsClient.GetTransactions(context.Background(), transactionsRequest)

	var transactions []models.TransactionResponse

	if err != nil {
		return transactions, false, err
	}

	for _, ts := range response.Transactions {
		loopTS := &models.TransactionResponse{
			TsID:      ts.TsID,
			FromUser:  ts.FromID,
			FromName:  ts.FromName,
			ToUser:    ts.ToID,
			ToName:    ts.ToName,
			Amount:    ts.Amount,
			Message:   ts.Message,
			CreatedAt: ts.CreateAt,
		}
		transactions = append(transactions, *loopTS)
	}

	return transactions, true, nil
}
