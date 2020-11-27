package grpc

import (
	"log"
	"time"

	internalDB "github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/storage"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

//Server Movement
type Server struct{}

var db *gorm.DB

//StartDB start the db for gRPC
func StartDB() {
	db = internalDB.ConnectDB()

	if db == nil {
		log.Fatalln("DB no conneted")
	}
}

//CloseDB close the db for gRPC
func CloseDB() {
	db.Close()
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

	DBService := storage.NewMovementStorageService(db)

	success, err := DBService.NewMovement(&newMove)

	if err != nil {
		return &MovementResponse{}, err
	}

	return &MovementResponse{Sucess: success}, nil
}
