package rmq

import (
	"context"
	"fmt"
	pr_result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-controller/model"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(ctx context.Context, sub chan<- pr_result.Event, ID model.MsgID)
	Close() error
}

type rmqConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	query   string
	router  Router
}

func NewConsumer(connection *amqp.Connection, queryName string) (Consumer, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(queryName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error queue connection %s: %w", queryName, err)
	}

	results, err := channel.Consume(
		q.Name,
		"calc-controller",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	toRouter := make(chan MessageFromRMQ)
	go func() {
		for event := range results {
			msg := MessageFromRMQ{
				ContentType: event.ContentType,
				Timestamp:   event.Timestamp,
				MessageID:   model.MsgID(event.MessageId),
				Body:        event.Body,
			}
			toRouter <- msg
		}
		close(toRouter)
	}()

	return &rmqConsumer{
		conn:    connection,
		channel: channel,
		query:   queryName,
		router:  InitFilter(toRouter),
	}, nil
}

func (c *rmqConsumer) Consume(ctx context.Context, sub chan<- pr_result.Event, ID model.MsgID) {
	c.router.AddRoute(ID, sub)
}

func (c *rmqConsumer) Close() error {
	return c.channel.Close()
}
