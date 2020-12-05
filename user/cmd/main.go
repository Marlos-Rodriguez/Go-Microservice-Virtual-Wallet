package main

import (
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/api"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/grpc"
)

func main() {
	go grpc.NewGRPCServer()
	api.Start()
}
