package main

import (
	"fmt"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/parser"
	"time"

	schema "gAD-System/internal/proto/expression/event"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println("Error logger sync")
		}
	}()

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Error("failed to init cfg from with envconfig")
	}

	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://%s", cfg.RMQConfig.Server))
	if err != nil {
		logger.Fatal("failed to connect to rabbitmq:", zap.Error(err))
	}
	defer rmqConn.Close()

	ch, err := rmqConn.Channel()
	if err != nil {
		logger.Error("failed to open channel")
		return
	}
	defer ch.Close()

	exprs, err := ch.Consume(
		cfg.RMQConfig.PubQueryName,
		"calc-worker",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("failed to open consume")
		return
	}

	end := make(chan struct{})

	go func() {
		for msg := range exprs {
			expr, err := protoToMsg(msg.Body)
			if err != nil {
				logger.Error("proto to msg error", zap.Error(err), zap.String("msg id", msg.MessageId))
			}

			result, err := parser.CalculateSimpleExpression(expr)
			if err != nil {
				logger.Error("err in calc expr", zap.Error(err), zap.String("msg id", msg.MessageId))
			}

			body, err := msgToProtoBytes(result)
			if err != nil {
				logger.Error("msg to proto error", zap.Error(err), zap.String("msg id", msg.MessageId))
			}

			err = ch.Publish("", cfg.RMQConfig.SubQueryName, false, false, amqp.Publishing{
				ContentType: msg.ContentType,
				MessageId:   msg.MessageId,
				Timestamp:   time.Now(),
				Body:        body,
			})
			if err != nil {
				logger.Error("msg rmq publishing", zap.Error(err), zap.String("msg id", msg.MessageId), zap.String("queue name", cfg.RMQConfig.SubQueryName))
			}
		}
		end <- struct{}{}
	}()

	<-end
}

func msgToProtoBytes(message string) ([]byte, error) {
	event := schema.Event{Expression: message}
	out, err := proto.Marshal(&event)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func protoToMsg(data []byte) (string, error) {
	event := schema.Event{}
	if err := proto.Unmarshal(data, &event); err != nil {
		return "", err
	}
	return event.Expression, nil
}
