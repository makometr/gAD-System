package main

import (
	"fmt"
	"gAD-System/services/gad-manager/config"
	"gAD-System/services/gad-manager/domain"
	"gAD-System/services/gad-manager/server"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println("Error logger sync")
		}
	}()

	cfg := config.InitConfig()

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

	if err := server.InitREST(cfg, &handlers); err != nil {
		logger.Error("failed to init GIN-REST API server", zap.Error(err))
	}

}
