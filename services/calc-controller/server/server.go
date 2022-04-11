package server

import (
	"context"
	"fmt"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	pr_result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-controller/parser"
	"gAD-System/services/calc-controller/rmq"
	"strconv"
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
	exprsToCalculate := request.GetExpression()
	results := make([]string, len(exprsToCalculate))

	// отравляем их "асинхронно" на вычисления
	var wg sync.WaitGroup
	wg.Add(len(exprsToCalculate))
	for i, expr := range exprsToCalculate {
		go func(gIndex int, expr string) {
			parsedExpr, err := parser.ParseBinaryExpression(expr)
			if err != nil {
				fmt.Println("Error in ParseBinaryExpression():", err)
			}

			result, err := s.exprCalculator.CalculateExpression(*parsedExpr)
			if err != nil {
				fmt.Println("Error in DoCalculate():", err)
			}

			var finalResult string
			switch res := result.Result.(type) {
			case *pr_result.Event_Product:
				finalResult = strconv.FormatInt(res.Product, 10)
			case *pr_result.Event_ErrorMsg:
				finalResult = "error: " + res.ErrorMsg
			default:
				finalResult = "internal error"
				fmt.Println("enexpected value from protobuf conversation")
			}

			results[gIndex] = finalResult
			wg.Done()
		}(i, expr)
	}

	wg.Wait()

	reply := pb.CalculatorReply{Result: results}
	return &reply, nil
}
