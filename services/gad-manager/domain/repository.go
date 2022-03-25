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
	doCalculate([]string) ([]string, error)
}

type CalcRepo struct {
	calcClient pb.CalculatorServiceClient
}

func NewCalcRepository(conn *grpc.ClientConn) Repository {
	c := pb.NewCalculatorServiceClient(conn)
	return &CalcRepo{calcClient: c}
}

func (cr CalcRepo) doCalculate(exprs []string) ([]string, error) {
	ctx := context.TODO()
	r, err := cr.calcClient.DoCalculate(ctx, &pb.CalculatorRequest{Expression: exprs})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errCalcRPCResponse, err)
	}

	return r.Result, nil
}
