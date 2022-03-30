package rmq

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/streadway/amqp"
)

var ErrProtobuffSerialize = errors.New("error converting message to proto bytes")
var ErrAMQSend = errors.New("error seding msg to rmq")

type Producer interface {
	SendExpresion(ctx context.Context, expr ExpressionWithID) error
	Close() error
}

type rmqProducer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	query string
	// in    chan Message
}

func NewProducer(connection *amqp.Connection, queryName string) (Producer, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	producer := &rmqProducer{
		conn:  connection,
		ch:    channel,
		query: queryName,
	}

	return producer, nil
}

func (p *rmqProducer) SendExpresion(ctx context.Context, expr ExpressionWithID) error {
	serialize, err := MsgToProtoBytes(expr.Expr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProtobuffSerialize, err)
	}

	err = p.ch.Publish("", p.query, false, false, amqp.Publishing{
		ContentType: "text/plain",
		MessageId:   string(expr.Id),
		Timestamp:   time.Now(),
		Body:        serialize,
	})

	if err != nil {
		return fmt.Errorf("%w: %v", ErrAMQSend, err)
	}

	return nil
}

func (p *rmqProducer) Close() error {
	return p.ch.Close()
}
