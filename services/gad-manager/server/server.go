package server

import (
	"fmt"
	"gAD-System/services/gad-manager/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitREST inits REST API server for outer clients
func InitREST(cfg *config.Config, h *Handlers) error {
	r := newRouter(h)
	if err := r.Run(fmt.Sprintf(":%s", cfg.GMConfig.Port)); err != nil {
		return err
	}
	return nil
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
