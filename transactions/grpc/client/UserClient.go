package client

import (
	"context"
	"errors"
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	"google.golang.org/grpc"
)

var userClient models.UserServiceClient
var userConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startUserClient() {
	urlTarget := environment.AccessENV("USER_GRPC")

	if urlTarget == "" {
		log.Fatalln("Error in Access to User_GRPC URL in Transaction Service")
	}
	userConn, err := grpc.Dial(urlTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	userClient = models.NewUserServiceClient(userConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeUserClient() {
	userConn.Close()
}

//CheckUserTransaction Create a new movement in DB
func CheckUserTransaction(fromID string, toID string, amount float32, password string) (bool, error) {
	newTransaction := &models.CheckTransactionRequest{
		FromID:   fromID,
		ToID:     toID,
		Amount:   float64(amount),
		Password: password,
	}

	response, err := userClient.CheckUsersTransactions(context.Background(), newTransaction)

	if err != nil {
		return false, err
	}

	if !response.Actives || !response.Exits || !response.Enough {
		return false, errors.New("Transaction is not posible")
	}

	return true, nil
}

//MakeTransaction Between Users
func MakeTransaction(fromID string, toID string, amount float32) (bool, error) {
	newTransaction := &models.TransactionRequest{
		FromID: fromID,
		ToID:   toID,
		Amount: float64(amount),
	}

	response, err := userClient.MakeTransaction(context.Background(), newTransaction)

	if err != nil {
		return false, err
	}

	return response.Sucess, nil
}
