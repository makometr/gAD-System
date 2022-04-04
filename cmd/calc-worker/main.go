package main

import (
	"fmt"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/rmq"
	"os"
	"os/signal"
	"sync"
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
	logger.Info("RMQ connection succeeded:", zap.Any("rmq", cfg.RMQConfig))

	errorChan := make(chan error)
	go func() {
		for err := range errorChan {
			logger.Error("error in calculate chain:", zap.Error(err))
		}
	}()

	workerCount := 8
	wg := sync.WaitGroup{}
	wg.Add(workerCount)
	// for i := 0; i < workerCount; i++ {
	// 	go func() {
	// 		rmqConn.CalculateExpressions(errorChan, parser.CalculateSimpleExpression)
	// 		wg.Done()
	// 	}()
	// }

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	fmt.Print("RMQ connection closing...")
	if err := rmqConn.Close(); err != nil {
		logger.Error("RMQ conn close error:", zap.Error(err))
	}
	fmt.Print(" closed.\n")

	fmt.Print("Workers are closing...")
	wg.Wait()
	fmt.Print(" closed.\n")

	close(errorChan)
}
