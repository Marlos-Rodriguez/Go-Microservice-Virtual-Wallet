package main

import (
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/server"
)

func main() {
	log.Println("Start Movement gRPC Server")

	server.NewGRPCServer()
}
