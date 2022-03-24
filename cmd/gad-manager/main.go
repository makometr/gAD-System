package main

import (
	"gAD-System/services/gad-manager/config"
	"gAD-System/services/gad-manager/domain"
	"gAD-System/services/gad-manager/server"
)

func main() {
	cfg := config.InitConfig()

	calcConn := server.InitCalculateRPC(cfg)
	defer calcConn.Close()
	calcRepo := domain.NewCalcRepository(calcConn)
	calculator := domain.NewCalculator(calcRepo)
	handlers := server.Handlers{Calculator: calculator}

	_ = server.InitREST(cfg, &handlers)
}
