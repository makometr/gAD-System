package server

import (
	"context"
	pb "gAD-System/internal/proto/grpc/calculator/service"
)

type calculatorServer struct {
	pb.CalculatorServiceServer
}

func NewCalculatorServer() *calculatorServer {
	return &calculatorServer{}
}

func (s *calculatorServer) DoCalculate(ctx context.Context, request *pb.CalculatorRequest) (*pb.CalculatorReply, error) {
	req := request.GetExpression()

	res := make([]string, 0, len(req))
	for _, item := range req {
		res = append(res, item+"!")
	}

	reply := pb.CalculatorReply{Result: res}
	return &reply, nil
}
