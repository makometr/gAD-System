package server

import (
	"fmt"
	"gAD-System/services/gad-manager/config"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LaunchREST inits REST API server for outer clients
func LaunchREST(cfg *config.Config, h *Handlers) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.GMConfig.Port),
		Handler: newRouter(h),
	}
}

// InitCalculateRPC provide grps connection
func InitCalculateRPC(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.CCConfig.Server, cfg.CCConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
