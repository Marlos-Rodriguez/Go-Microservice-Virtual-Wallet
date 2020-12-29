package main

import (
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/api"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/grpc"
)

func main() {
	go grpc.NewGRPCServer()
	api.Start()
}
