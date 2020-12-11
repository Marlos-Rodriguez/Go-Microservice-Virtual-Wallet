package grpc

import (
	"log"
	"net"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/grpc/server"

	"google.golang.org/grpc"
)

//NewGRPCServer Create new gRPC server
func NewGRPCServer() {
	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("Failed to listen on Port :9001 %v", err)
	}

	s := server.Server{}

	grpcServer := grpc.NewServer()

	server.RegisterAuthServiceServer(grpcServer, &s)

	server.GetStorageService()
	defer server.CloseDB()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on Port :9001 %v", err)
	}
}
