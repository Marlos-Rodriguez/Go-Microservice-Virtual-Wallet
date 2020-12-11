package grpc

import (
	"log"
	"time"

	internalDB "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/storage"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

//Server Movement
type Server struct{}

var dbService *storage.MovementService

//StartDB start the db for gRPC
func StartDB() {
	db := internalDB.ConnectDB("MOVE")

	if db == nil {
		log.Fatalln("DB no conneted")
	}

	dbService = storage.NewMovementStorageService(db)
}

//CloseDB close the db for gRPC
func CloseDB() {
	dbService.CloseDB()
}

//CreateMovement Create a New movement Server method
func (s *Server) CreateMovement(ctx context.Context, move *MovementRequest) (*MovementResponse, error) {
	newMove := *&models.Movement{
		MovementID: uuid.New(),
		Relation:   move.Relation,
		Change:     move.Change,
		Origin:     move.Origin,
		CreatedAt:  time.Now(),
	}

	success, err := dbService.NewMovement(&newMove)

	if err != nil {
		return &MovementResponse{}, err
	}

	return &MovementResponse{Sucess: success}, nil
}
