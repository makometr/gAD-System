package rmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(ctx context.Context, sub chan<- ExpressionWithID, ID MsgID)
}

type rmqConsumer struct {
	conn   *amqp.Connection
	query  string
	router Router
}

func NewConsumer(connection *amqp.Connection, queryName string) (Consumer, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	results, err := channel.Consume(
		queryName,
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
				MessageID:   MsgID(event.MessageId),
				Body:        event.Body,
			}
			fmt.Println("readed ig loop in consumer:", msg.MessageID, string(msg.Body))
			ch <- msg
		}
		close(ch)
	}()

	return &rmqConsumer{
		conn:   connection,
		query:  queryName,
		router: InitFilter(ch),
	}, nil
}

func (c *rmqConsumer) Consume(ctx context.Context, sub chan<- ExpressionWithID, ID MsgID) {
	c.router.AddRoute(ID, sub)
}
