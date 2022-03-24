package domain

import (
	"google.golang.org/grpc"
)

type Repository interface {
	doCalculate(string) (string, error)
}

type CalcRepo struct {
	Conn *grpc.ClientConn
}

func NewCalcRepository(conn *grpc.ClientConn) Repository {
	return &CalcRepo{Conn: conn}
}

func (r CalcRepo) doCalculate(string) (string, error) {
	// TODO grpc-client
	return "100", nil
}
