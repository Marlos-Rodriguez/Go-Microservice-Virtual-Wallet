package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/transactions/grpc/server"
)

//NewGRPCServer Create new gRPC server
func NewGRPCServer() {
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Failed to listen on Port :9003 %v", err)
	}

	s := server.Server{}

	grpcServer := grpc.NewServer()

	server.RegisterTransactionServiceServer(grpcServer, &s)

	server.GetStorageService()
	defer server.CloseDB()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on Port :9001 %v", err)
	}
}
