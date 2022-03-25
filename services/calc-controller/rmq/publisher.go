package rmq

import (
	"context"
	"github.com/streadway/amqp"
	"time"
)

type Message struct {
	ContentType string
	Timestamp   time.Time
	MessageID   string
	Body        []byte
}

type Publisher interface {
	Publish(ctx context.Context, message Message) error
}

type publisherImpl struct {
	channel   *amqp.Channel
	queryName string
}

func NewPublisher(channel *amqp.Channel, queryName string) Publisher {
	return &publisherImpl{
		channel:   channel,
		queryName: queryName,
	}
}

func (p *publisherImpl) Publish(ctx context.Context, message Message) error {
	err := p.channel.Publish(
		"",
		p.queryName,
		false,
		false,
		amqp.Publishing{
			ContentType: message.ContentType,
			MessageId:   message.MessageID,
			Timestamp:   message.Timestamp,
			Body:        message.Body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
