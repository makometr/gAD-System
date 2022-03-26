package server

import (
	"context"
	"fmt"
	schema "gAD-System/internal/proto/expression/event"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/rmq"
	"google.golang.org/protobuf/proto"
	"sync"
)

type calculatorServer struct {
	pb.CalculatorServiceServer
	publisher rmq.Publisher
	consumer  rmq.Consumer
}

func NewCalculatorServer(publisher rmq.Publisher, consumer rmq.Consumer) *calculatorServer {
	return &calculatorServer{
		publisher: publisher,
		consumer:  consumer,
	}
}

func (s *calculatorServer) DoCalculate(ctx context.Context, request *pb.CalculatorRequest) (*pb.CalculatorReply, error) {
	input := make(chan rmq.Message)
	output := make(chan rmq.Message)

	var wg sync.WaitGroup
	wg.Add(3)

	var results []string

	go func() {
		defer wg.Done()
		payload := request.GetExpression()
		for _, event := range payload {
			serialize, err := msgToProtoBytes(event)
			if err != nil {
				fmt.Printf("error converting message to proto bytes: %v", err)
				continue
			}
			msg := rmq.Message{
				ContentType: "text/plain",
				MessageID:   event,
				Body:        serialize,
			}
			input <- msg
		}
		close(input)
	}()

	go func() {
		defer wg.Done()
		err := s.publisher.Publish(ctx, input)
		if err != nil {
			fmt.Printf("error while publishing events: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := s.consumer.Consume(ctx, output)
		if err != nil {
			fmt.Printf("error while consuming events: %v", err)
		}
	}()

	wg.Wait()
	for event := range output {
		fmt.Printf("New event in DoCalculation: %s", event)
		results = append(results, event.MessageID)
	}

	reply := pb.CalculatorReply{Result: append(results, "EVERYTHING WROCKS!")}
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
