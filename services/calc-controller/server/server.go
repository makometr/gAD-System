package server

import (
	"context"
	"fmt"
	expr "gAD-System/internal/proto/expression/event"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-controller/rmq"
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

	test_expr := expr.Event{Lhs: 100, Rhs: 200, Operation: expr.Operation_PLUS}

	// отравляем их "асинхронно" на вычисления
	var results []string

	var test_result result.Event
	test_result, err := s.exprCalculator.CalculateExpression(test_expr)
	if err != nil {
		fmt.Println("Error in DoCalculate():", err)
	}
	results = append(results, test_result.String())

	// var wg sync.WaitGroup
	// wg.Add(len(payload))
	// for _, expr := range payload {
	// 	go func(expr string) {
	// 		result, err := s.exprCalculator.CalculateExpression(model.Expression{Lhs: 100, Rhs: 200, Oper: model.Plus})
	// 		if err != nil {
	// 			fmt.Println("Error in DoCalculate():", err)
	// 		}
	// 		results = append(results, strconv.Itoa(int(result.Result)))
	// 		wg.Done()
	// 	}(expr)
	// }

	// wg.Wait()

	reply := pb.CalculatorReply{Result: append(results, "EVERYTHING WROCKS!")}
	return &reply, nil
}
