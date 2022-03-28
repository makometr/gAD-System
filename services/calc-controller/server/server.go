package server

import (
	"context"
	"fmt"
	schema "gAD-System/internal/proto/expression/event"
	pb "gAD-System/internal/proto/grpc/calculator/service"
	"gAD-System/services/calc-controller/rmq"
	"google.golang.org/protobuf/proto"
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
	payload := request.GetExpression()

	input := make(chan rmq.Message, len(payload))
	output := make(chan rmq.Message)
	err := make(chan error)

	var results []string

	go func() {
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
		fmt.Println("go1 finished")
	}()

	go func() {
		err := s.publisher.Publish(ctx, input)
		if err != nil {
			fmt.Printf("error while publishing events: %v", err)
		}
		fmt.Println("go2 finished")
	}()

	go func() {
		err := s.consumer.Consume(ctx, output)
		if err != nil {
			fmt.Printf("error while consuming events: %v", err)
		}
		fmt.Println("go3 finished")
	}()

	go func() {
		for i := 0; i < len(payload); i++ {
			event := <-output
			fmt.Printf("New event in DoCalculation: %s", event)
			results = append(results, event.MessageID)
		}
		err <- nil
	}()

	<-err
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
