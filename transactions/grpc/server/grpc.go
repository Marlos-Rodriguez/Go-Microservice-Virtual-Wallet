package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

//NewGRPCServer Create new gRPC server
func NewGRPCServer() {
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Failed to listen on Port :9003 %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	RegisterTransactionServiceServer(grpcServer, &s)

	GetStorageService()
	defer CloseDB()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on Port :9001 %v", err)
	}
}
