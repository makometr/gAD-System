package domain

import (
	"google.golang.org/grpc"
)

type Repository interface {
	doCalculate(string) (string, error)
	Close() error
}

type CalcRepo struct {
	Conn *grpc.ClientConn
}

func NewCalcRepository(conn *grpc.ClientConn) Repository {
	return &CalcRepo{Conn: conn}
}

func (r CalcRepo) doCalculate(string) (string, error) {
	// r.Conn.
	// TODO
	return "100", nil
}

func (r CalcRepo) Close() error {
	return r.Conn.Close()
}
