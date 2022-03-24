package main

import (
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("could not start listen on 8081: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer())
	grpcServer.Serve(listen)
}
