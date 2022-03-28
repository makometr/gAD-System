package rmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, pub <-chan Message) error
}

type rmqPublisher struct {
	conn  *amqp.Connection
	query string
}

func NewPublisher(connection *amqp.Connection, queryName string) Publisher {
	return &rmqPublisher{
		conn:  connection,
		query: queryName,
	}
}

func (p *rmqPublisher) Publish(ctx context.Context, pub <-chan Message) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}
	fmt.Println("channel created")
	defer channel.Close()

	quit := make(chan error)
	go func() {
		for event := range pub {
			err = channel.Publish("", p.query, false, false, amqp.Publishing{
				ContentType: event.ContentType,
				MessageId:   event.MessageID,
				Timestamp:   time.Now(),
				Body:        event.Body,
			})
			if err != nil {
				quit <- err
			}
			fmt.Println("published", event)
		}
		quit <- nil
	}()

	if err = <-quit; err != nil {
		return err
	}
	fmt.Println("Publishing finished")
	return nil
}
