package client

import (
	"context"
	"log"

	movements "github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/grpc"
	"google.golang.org/grpc"
)

var moveClient movements.MovementServiceClient
var moveConn *grpc.ClientConn

//StartMoveClient Start the client for movement gRPC
func startMoveClient() {
	moveConn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	moveClient = movements.NewMovementServiceClient(moveConn)
}

//CloseMoveClient Close the client for movement gRPC
func closeMoveClient() {
	moveConn.Close()
}

//CreateMovement Create a new movement in DB
func CreateMovement(relation string, change string, origin string) (bool, error) {
	newMovement := &movements.MovementRequest{
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
