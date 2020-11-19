package server

import (
	"log"
	"net"

	MoveServer "github.com/Marlos-Rodriguez/go-postgres-wallet-back/movements/grpc"
	"google.golang.org/grpc"
)

//NewGRPCServer Create new gRPC server
func NewGRPCServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on Port :9000 %v", err)
	}

	s := MoveServer.Server{}

	grpcServer := grpc.NewServer()

	MoveServer.RegisterMovementServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on Port :9000 %v", err)
	}
}
