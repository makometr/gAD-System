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

	ch := make(chan Message)
	go func() {
		fmt.Println("Waiting to consume...")
		for event := range results {
			fmt.Printf("New event is coming: %s", event.MessageId)
			msg := Message{
				ContentType: event.ContentType,
				Timestamp:   event.Timestamp,
				MessageID:   MsgID(event.MessageId),
				Body:        event.Body,
			}
			ch <- msg
		}
		close(ch)
		fmt.Println("Consuming finished")
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
