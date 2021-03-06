package main

import (
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/api"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc"
)

func main() {
	log.Println("Start Auth gRPC")
	go grpc.NewGRPCServer()
	log.Println("Start Auth Server")
	api.Start()
}
