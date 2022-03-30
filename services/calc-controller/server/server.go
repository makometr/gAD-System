package server

import (
	"context"
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/rmq"
	"sync"
)

type calculatorServer struct {
	pb.CalculatorServiceServer
	exprCalculator rmq.ExprCalculator
}

func NewCalculatorServer(exprCalc rmq.ExprCalculator) *calculatorServer {
	return &calculatorServer{
		exprCalculator: exprCalc,
	}
}

func (s *calculatorServer) DoCalculate(ctx context.Context, request *pb.CalculatorRequest) (*pb.CalculatorReply, error) {
	payload := request.GetExpression()

	// отравляем их "асинхронно" на вычисления
	var wg sync.WaitGroup
	var results []string
	wg.Add(len(payload))
	for _, expr := range payload {
		go func(expr string) {
			result, err := s.exprCalculator.CalculateExpression(expr)
			if err != nil {
				fmt.Println("Error in DoCalculate():", err)
			}
			results = append(results, result)
			wg.Done()
		}(expr)
	}

	wg.Wait()

	reply := pb.CalculatorReply{Result: append(results, "EVERYTHING WROCKS!")}
	return &reply, nil
}
