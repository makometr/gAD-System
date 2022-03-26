package main

import (
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/config"
	"gAD-System/services/calc-controller/rmq"
	"gAD-System/services/calc-controller/server"

	"github.com/streadway/amqp"

	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println("Error logger sync")
		}
	}()

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Error("failed to init cfg from with envconfig")
	}

	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%s/", cfg.RMQCalc.Port))
	if err != nil {
		logger.Fatal("failed to connect to rabbitmq:", zap.Error(err))
	}
	defer rmqConn.Close()

	rmqPub := rmq.NewPublisher(rmqConn, cfg.RMQCalc.PubQName)
	rmqSub := rmq.NewConsumer(rmqConn, cfg.RMQCalc.SubQName)

	listen, err := net.Listen("tcp", cfg.RPCCalc.Port)
	if err != nil {
		logger.Error("failed to init RPC connection:", zap.Error(err))
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer(rmqPub, rmqSub))
	grpcServer.Serve(listen)
}
