package server

import (
	"context"
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/rmq"
	"sync"
	"time"
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
	// payload := request.GetExpression()

	// парсим дерево и получаем два выражения
	// lhs := rmq.ExpressionWithID{Expr: "100+100", Id: "11"}
	// rhs := rmq.ExpressionWithID{Expr: "2+2", Id: "22"}

	// отравляем их "асинхронно" на вычисления
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		time.Sleep(1 * time.Second)
		result, err := s.exprCalculator.CalculateExpression("100+100")
		fmt.Println("Recieve result:", result, err)
		wg.Done()
	}()

	// отравляем их "асинхронно" на вычисления
	go func() {
		result, err := s.exprCalculator.CalculateExpression("2+2")
		fmt.Println("Recieve result:", result, err)
		wg.Done()
	}()

	wg.Wait()

	var results []string
	reply := pb.CalculatorReply{Result: append(results, "EVERYTHING WROCKS!")}
	return &reply, nil
}
