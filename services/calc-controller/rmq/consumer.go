package rmq

import (
	"context"
	"fmt"
	"gAD-System/services/calc-controller/model"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(ctx context.Context, sub chan<- model.ResultFromCalc, ID model.MsgID)
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

	fmt.Println("Consumer listened:", queryName)

	ch := make(chan Message)
	go func() {
		for event := range results {
			msg := Message{
				ContentType: event.ContentType,
				Timestamp:   event.Timestamp,
				MessageID:   model.MsgID(event.MessageId),
				Body:        event.Body,
			}
			fmt.Println("readed ig loop in consumer:", msg.MessageID, string(msg.Body))
			ch <- msg
		}
		close(ch)
	}()

	return &rmqConsumer{
		conn:    connection,
		channel: channel,
		query:   queryName,
		router:  InitFilter(ch),
	}, nil
}

func (c *rmqConsumer) Consume(ctx context.Context, sub chan<- model.ResultFromCalc, ID model.MsgID) {
	c.router.AddRoute(ID, sub)
}

func (c *rmqConsumer) Close() error {
	return c.channel.Close()
}
