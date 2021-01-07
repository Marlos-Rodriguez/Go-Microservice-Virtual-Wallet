package client

import (
	"context"
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
	"google.golang.org/grpc"
)

var authClient models.AuthServiceClient
var authConn *grpc.ClientConn

//StartAuthClient Start the client for Auth gRPC
func startAuthClient() {
	urlTarget := environment.AccessENV("AUTH_GRPC")

	if urlTarget == "" {
		log.Fatalln("Error in Access to GRPC URL in User Service")
	}
	authConn, err := grpc.Dial(urlTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	authClient = models.NewAuthServiceClient(authConn)
}

//CloseAuthClient Close the client for movement gRPC
func closeAuthClient() {
	authConn.Close()
}

//UpdateAuthCache Update the User username or email
func UpdateAuthCache(oldUsername string, newUsername string, oldEmail string, newEmail string) (bool, error) {
	User := &models.NewUserInfo{
		OldUsername: oldUsername,
		NewUsername: newUsername,
		OldEmail:    oldEmail,
		NewEmail:    newEmail,
	}

	response, err := authClient.ChangeAuthCache(context.Background(), User)

	if err != nil {
		return false, err
	}

	return response.Success, nil
}
