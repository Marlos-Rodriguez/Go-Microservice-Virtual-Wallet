package api

import (
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	"github.com/jinzhu/gorm"
)

//Start Start a new User server API
func Start() {
	var DB *gorm.DB

	DB = storage.ConnectDB()
	defer DB.Close()

	app := routes(DB)
	createServer(app)
}
