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
	ch    *amqp.Channel
	query string
}

func NewPublisher(connection *amqp.Connection, queryName string) (Publisher, *amqp.Channel, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	return &rmqPublisher{
		conn:  connection,
		ch:    channel,
		query: queryName,
	}, channel, nil
}

func (p *rmqPublisher) Publish(ctx context.Context, pub <-chan Message) error {
	quit := make(chan error)
	go func() {
		for event := range pub {
			err := p.ch.Publish("", p.query, false, false, amqp.Publishing{
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

	if err := <-quit; err != nil {
		return err
	}
	fmt.Println("Publishing finished")
	return nil
}
