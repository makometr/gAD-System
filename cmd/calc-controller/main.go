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
		return
	}

	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://%s", cfg.RMQConfig.Server))
	if err != nil {
		logger.Fatal("failed to connect to rabbitmq:", zap.Error(err))
		return
	}
	defer rmqConn.Close()

	rmqPub, err := rmq.NewProducer(rmqConn, cfg.RMQConfig.PubQueryName)
	if err != nil {
		logger.Fatal("failed to create new publisher", zap.Error(err))
		return
	}
	defer rmqPub.Close()

	rmqSub, err := rmq.NewConsumer(rmqConn, cfg.RMQConfig.SubQueryName)
	if err != nil {
		logger.Fatal("failed to create new consumer", zap.Error(err))
		return
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.CCConfig.Port))
	if err != nil {
		logger.Error("failed to init RPC connection:", zap.Error(err))
		return
	}

	exprCalculator := rmq.NewRemoteCalculator(rmqPub, rmqSub)

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer(exprCalculator))
	grpcServer.Serve(listen)
}
