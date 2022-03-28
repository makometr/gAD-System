package rmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(ctx context.Context, sub chan<- Message) error
}

type rmqConsumer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	query string
}

func NewConsumer(connection *amqp.Connection, queryName string) (Consumer, *amqp.Channel, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	return &rmqConsumer{
		conn:  connection,
		query: queryName,
	}, channel, nil
}

func (c *rmqConsumer) Consume(ctx context.Context, sub chan<- Message) error {
	results, err := c.ch.Consume(
		c.query,
		"calc-controller",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	quit := make(chan error)

	go func() {
		fmt.Println("Waiting to consume...")
		for event := range results {
			fmt.Printf("New event is comming: %s", event.MessageId)
			msg := Message{
				ContentType: event.ContentType,
				Timestamp:   event.Timestamp,
				MessageID:   event.MessageId,
				Body:        event.Body,
			}
			sub <- msg
		}
		close(sub)
		quit <- nil
	}()

	if err = <-quit; err != nil {
		return err
	}
	fmt.Println("Consuming finished")
	return nil
}
