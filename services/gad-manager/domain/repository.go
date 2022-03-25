package domain

import (
	"context"
	"errors"
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"

	"google.golang.org/grpc"
)

var errCalcRPCResponse = errors.New("receive calc result grpc failure")

type Repository interface {
	DoCalculate([]string) ([]string, error)
}

type сalcRepo struct {
	calcClient pb.CalculatorServiceClient
}

func NewCalcRepository(conn *grpc.ClientConn) Repository {
	c := pb.NewCalculatorServiceClient(conn)
	return &сalcRepo{calcClient: c}
}

func (cr сalcRepo) DoCalculate(exprs []string) ([]string, error) {
	ctx := context.TODO()
	r, err := cr.calcClient.DoCalculate(ctx, &pb.CalculatorRequest{Expression: exprs})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errCalcRPCResponse, err)
	}

	return r.Result, nil
}
