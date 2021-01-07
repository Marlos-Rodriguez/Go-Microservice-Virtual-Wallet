package client

import (
	"context"
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/models"
	"google.golang.org/grpc"
)

var moveClient models.MovementServiceClient
var moveConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startMoveClient() {
	urlTarget := environment.AccessENV("MOVEMENT_GRPC")

	if urlTarget == "" {
		log.Fatalln("Error in Access to MOVEMENT_GRPC URL in Transaction Service")
	}
	moveConn, err := grpc.Dial(urlTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	moveClient = models.NewMovementServiceClient(moveConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeMoveClient() {
	moveConn.Close()
}

//CreateMovement Create a new movement in DB
func CreateMovement(relation string, change string, origin string) (bool, error) {
	newMovement := &models.MovementRequest{
		Relation: relation,
		Change:   change,
		Origin:   origin,
	}

	response, err := moveClient.CreateMovement(context.Background(), newMovement)

	if err != nil {
		return false, err
	}

	return response.Sucess, nil
}
