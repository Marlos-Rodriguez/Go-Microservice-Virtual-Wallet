package storage

import (
	"context"
	"log"

	TSModels "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	TSserver "github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/server"
	"google.golang.org/grpc"
)

var tsClient TSserver.TransactionServiceClient
var tsConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startTransactionClient() {
	moveConn, err := grpc.Dial(":9003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	tsClient = TSserver.NewTransactionServiceClient(moveConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeTransactionClient() {
	tsConn.Close()
}

//GetTransactions of user
func GetTransactions(id string) ([]TSModels.TransactionResponse, bool, error) {
	transactionsRequest := &TSserver.TransactionRequest{
		ID: id,
	}

	response, err := tsClient.GetTransactions(context.Background(), transactionsRequest)

	var transactions []TSModels.TransactionResponse

	if err != nil {
		return transactions, false, err
	}

	for _, ts := range response.Transactions {
		loopTS := &TSModels.TransactionResponse{
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
