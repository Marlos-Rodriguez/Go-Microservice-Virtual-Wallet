package grpc

import (
	"time"

	internalDB "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/storage"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

//Server Movement
type Server struct{}

//CreateMovement Create a New movement Server method
func (s *Server) CreateMovement(ctx context.Context, move *MovementRequest) (*MovementResponse, error) {
	newMove := *&models.Movement{
		MovementID: uuid.New(),
		Relation:   move.Relation,
		Change:     move.Change,
		Origin:     move.Origin,
		CreatedAt:  time.Now(),
	}

	newDB := internalDB.ConnectDB()

	DBService := storage.NewMovementStorageService(newDB)

	success, err := DBService.NewMovement(&newMove)

	if err != nil {
		return &MovementResponse{}, err
	}

	return &MovementResponse{Sucess: success}, nil
}
