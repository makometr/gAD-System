package server

import (
	"gAD-System/services/gad-manager/config"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalcRPCConn struct {
	Conn *grpc.ClientConn
}

func InitREST(cfg *config.Config) error {
	r := NewRouter()
	if err := r.Run(cfg.REST.Port); err != nil {
		return err
	}
	return nil
}

func InitCalculateRPC(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.Dial(cfg.RPCCalc.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
