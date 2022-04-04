package main

import (
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/config"
	"gAD-System/services/calc-controller/rmq"
	"gAD-System/services/calc-controller/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"

	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	defer log.Println("server stopped successful")
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
	defer func() {
		if err = rmqConn.Close(); err != nil {
			log.Printf("RMQ connection closed with error: %v\n", err)
			return
		}
		log.Println("RMQ connection closed.")
	}()

	rmqPub, err := rmq.NewProducer(rmqConn, cfg.RMQConfig.PubQueryName)
	if err != nil {
		logger.Fatal("failed to create new publisher", zap.Error(err))
	}
	defer func() {
		if err = rmqPub.Close(); err != nil {
			log.Printf("RMQ channel publish closed with error: %v\n", err)
			return
		}
		log.Println("RMQ channel publish closed.")
	}()

	rmqSub, err := rmq.NewConsumer(rmqConn, cfg.RMQConfig.SubQueryName)
	if err != nil {
		logger.Fatal("failed to create new consumer", zap.Error(err))
	}
	defer func() {
		if err = rmqSub.Close(); err != nil {
			log.Printf("RMQ channel subscriber closed with error: %v\n", err)
			return
		}
		log.Println("RMQ channel subscriber closed.")
	}()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.CCConfig.Port))
	if err != nil {
		logger.Error("failed to init RPC connection:", zap.Error(err))
		return
	}
	exprCalculator := rmq.NewRemoteCalculator(rmqPub, rmqSub)
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server.NewCalculatorServer(exprCalculator))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	grpcShutdowned := make(chan struct{})
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		grpcServer.GracefulStop()
		grpcShutdowned <- struct{}{}
	}()

	log.Println("starting grpc server")
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}

	<-grpcShutdowned
	log.Println("grpc-server shutdowned")

}
