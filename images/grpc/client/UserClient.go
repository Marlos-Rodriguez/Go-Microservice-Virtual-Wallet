package client

import (
	"context"
	"errors"
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/models/proto"
	"google.golang.org/grpc"
)

var userClient proto.UserServiceClient
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

	userClient = proto.NewUserServiceClient(userConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeUserClient() {
	userConn.Close()
}

//CheckUserTransaction Create a new movement in DB
func ChangeAvatar(url string, ID string) (bool, error) {
	newAvatar := &proto.AvatarRequest{
		ID:  ID,
		Url: url,
	}

	response, err := userClient.ChangeAvatar(context.Background(), newAvatar)

	if err != nil {
		return false, err
	}

	if !response.Sucess {
		return false, errors.New("Avatar no changed in user service")
	}

	return true, nil
}
