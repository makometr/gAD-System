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
	query string
}

func NewConsumer(connection *amqp.Connection, queryName string) Consumer {
	return &rmqConsumer{
		conn:  connection,
		query: queryName,
	}
}

func (c *rmqConsumer) Consume(ctx context.Context, sub chan<- Message) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	fmt.Println("Consume channel created")

	results, err := channel.Consume(
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
