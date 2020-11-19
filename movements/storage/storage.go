package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/models"
)

//MovementService struct
type MovementService struct {
	db *gorm.DB
}

//NewMovementStorageService Create a new DB movement service
func NewMovementStorageService(db *gorm.DB) *MovementService {

	db.AutoMigrate(&models.Movement{})

	return &MovementService{db: db}
}

//NewMovement Create a new movement
func (s *MovementService) NewMovement(m *models.Movement) (bool, error) {
	//Create new movement in DB
	if err := s.db.Create(&m).Error; err != nil {
		return false, nil
	}

	return true, nil
}
