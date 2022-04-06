package main

import (
	"context"
	"fmt"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/rmq"
	"gAD-System/services/calc-worker/worker"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
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

	rmqConn, err := rmq.InitConnection(cfg)
	if err != nil {
		logger.Error("rmq connection init error ", zap.Error(err))
		return
	}
	rmqInput, err := rmq.NewRMQInputStream(cfg, rmqConn)
	if err != nil {
		logger.Error("init input rmq channel erorr:", zap.Error(err))
		return
	}
	rmqOutput, err := rmq.NewRMQOutputStream(cfg, rmqConn)
	if err != nil {
		logger.Error("init ouptut rmq channel erorr:", zap.Error(err))
		return
	}
	logger.Info("RMQ connection succeeded:", zap.Any("rmq", cfg.RMQConfig))

	// errorChan := make(chan error)
	// go func() {
	// 	for err := range errorChan {
	// 		logger.Error("error in calculate chain:", zap.Error(err))
	// 	}
	// }()

	var input rmq.InputExprStream = rmqInput
	var output rmq.OutputExprStream = rmqOutput

	ctx := context.TODO()
	workersDone, err := worker.StartWorkers(ctx, cfg, input, output)
	if err != nil {
		logger.Error("error while start workers:", zap.Error(err))
		return
	}
	fmt.Println("Workers started!")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutting down from user/docker/os
	<-termChan

	fmt.Print("RMQ input channel closing... ")
	if err := rmqInput.Close(); err != nil {
		logger.Error("RMQ input channel close error:", zap.Error(err))
	}
	fmt.Print("closed.\n")

	fmt.Print("Workers are closing... ")
	<-workersDone
	fmt.Print("closed.\n")

	fmt.Print("RMQ output channel closing... ")
	if err := rmqOutput.Close(); err != nil {
		logger.Error("RMQ output channel close error:", zap.Error(err))
	}
	fmt.Print("closed.\n")

	fmt.Print("RMQ connection closing... ")
	if err := rmqConn.Close(); err != nil {
		logger.Error("RMQ connection close error:", zap.Error(err))
	}
	fmt.Print("closed.\n")

	// close(errorChan)
}
