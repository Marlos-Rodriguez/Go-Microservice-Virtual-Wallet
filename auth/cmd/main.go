package main

import (
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc"
)

func main() {
	log.Println("Start Auth gRPC")
	grpc.NewGRPCServer()
}
