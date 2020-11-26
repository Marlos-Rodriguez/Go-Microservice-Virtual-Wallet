package client

import (
	"context"
	"log"

	movements "github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/grpc"
	"google.golang.org/grpc"
)

//CreateMovement Create a new movement in DB
func CreateMovement(relation string, change string, origin string) (bool, error) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := movements.NewMovementServiceClient(conn)

	newMovement := &movements.MovementRequest{
		Relation: relation,
		Change:   change,
		Origin:   origin,
	}

	response, err := c.CreateMovement(context.Background(), newMovement)

	if err != nil {
		return false, err
	}

	return response.Sucess, nil
}
