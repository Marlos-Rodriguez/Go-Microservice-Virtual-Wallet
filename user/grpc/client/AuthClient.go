package client

import (
	"context"
	"log"

	auth "github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc/server"
	"google.golang.org/grpc"
)

var authClient auth.AuthServiceClient
var authConn *grpc.ClientConn

//StartAuthClient Start the client for Auth gRPC
func startAuthClient() {
	authConn, err := grpc.Dial(":9002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	authClient = auth.NewAuthServiceClient(authConn)
}

//CloseAuthClient Close the client for movement gRPC
func closeAuthClient() {
	authConn.Close()
}

//UpdateAuthCache Update the User username or email
func UpdateAuthCache(oldUsername string, newUsername string, oldEmail string, newEmail string) (bool, error) {
	User := &auth.UserRequest{}

	response, err := authClient.ChangeAuthCache(context.Background(), User)

	if err != nil {
		return false, err
	}

	return response.Success, nil
}
