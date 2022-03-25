package server

import (
	"context"
	schema "gAD-System/internal/proto/expression/event"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/rmq"
	"google.golang.org/protobuf/proto"
	"time"
)

type calculatorServer struct {
	pb.CalculatorServiceServer
	publisher rmq.Publisher
}

func NewCalculatorServer(publisher rmq.Publisher) *calculatorServer {
	return &calculatorServer{publisher: publisher}
}

func (s *calculatorServer) DoCalculate(ctx context.Context, request *pb.CalculatorRequest) (*pb.CalculatorReply, error) {
	payload := request.GetExpression()
	for _, expr := range payload {
		message, err := msgToProtoBytes(expr)
		if err != nil {
			return nil, err
		}

		err = s.publisher.Publish(ctx, rmq.Message{
			ContentType: "test/plain",
			Timestamp:   time.Now(),
			MessageID:   expr,
			Body:        message,
		})
	}

	reply := pb.CalculatorReply{Result: []string{"EVERYTHING WROKS!"}}
	return &reply, nil
}

func msgToProtoBytes(message string) ([]byte, error) {
	event := &schema.Event{Expression: message}
	out, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	return out, nil
}
