package main

import (
	"context"
	"fmt"
	"gAD-System/services/gad-manager/config"
	"gAD-System/services/gad-manager/domain"
	"gAD-System/services/gad-manager/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx := context.TODO()
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
	}

	calcConn, err := server.InitCalculateRPC(cfg)
	if err != nil {
		logger.Error("failed to init grpc:", zap.Error(err))
		return
	}
	defer func() {
		if err := calcConn.Close(); err != nil {
			logger.Error("failed to close grpc:", zap.Error(err))
		}
	}()

	calcRepo := domain.NewCalcRepository(calcConn)
	calculator := domain.NewCalculator(calcRepo)
	handlers := server.Handlers{Calculator: calculator}

	serverClosed := make(chan struct{})
	srv := server.LaunchREST(cfg, &handlers)
	go func() {
		err := srv.ListenAndServe()
		log.Println("server closed: ", err)
		serverClosed <- struct{}{}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("server shutdowned with err:", err)
		}
	}()

	<-serverClosed
}
