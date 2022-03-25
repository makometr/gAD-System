package main

import (
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("could not start listen on 50051: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer())
	grpcServer.Serve(listen)
}
