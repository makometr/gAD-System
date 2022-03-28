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

	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://%s", cfg.RMQConfig.Server))
	if err != nil {
		logger.Fatal("failed to connect to rabbitmq:", zap.Error(err))
	}
	defer rmqConn.Close()

	rmqPub, channelPub, err := rmq.NewPublisher(rmqConn, cfg.RMQConfig.PubQueryName)
	if err != nil {
		logger.Fatal("failed to create new publisher")
	}
	defer channelPub.Close()
	rmqSub, channelSub, err := rmq.NewConsumer(rmqConn, cfg.RMQConfig.SubQueryName)
	if err != nil {
		logger.Fatal("failed to create new consumer")
	}
	defer channelSub.Close()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.CCConfig.Port))
	if err != nil {
		logger.Error("failed to init RPC connection:", zap.Error(err))
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer(rmqPub, rmqSub))
	grpcServer.Serve(listen)
}
