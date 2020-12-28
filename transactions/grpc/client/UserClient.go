package client

import (
	"context"
	"errors"
	"log"

	userGRPC "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/grpc/server"
	"google.golang.org/grpc"
)

var userClient userGRPC.UserServiceClient
var userConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startUserClient() {
	userConn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	userClient = userGRPC.NewUserServiceClient(userConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeUserClient() {
	userConn.Close()
}

//CheckUserTransaction Create a new movement in DB
func CheckUserTransaction(fromID string, toID string, amount float32, password string) (bool, error) {
	newTransaction := &userGRPC.CheckTransactionRequest{
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
	newTransaction := &userGRPC.TransactionRequest{
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
