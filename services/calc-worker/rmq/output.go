package rmq

import (
	"fmt"
	"gAD-System/services/calc-worker/config"
	"time"

	"github.com/streadway/amqp"
)

type RMQOutputStream struct {
	qName      string
	channelOut *amqp.Channel
}

func NewRMQOutputStream(cfg *config.Config, rmq *Connection) (OutputExprStream, error) {
	rmqChanOut, err := rmq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel error: %w", err)
	}
	return &RMQOutputStream{channelOut: rmqChanOut, qName: cfg.RMQConfig.QNameResult}, nil
}

func (c *RMQOutputStream) Result(input <-chan string) (<-chan struct{}, error) {
	q, err := c.channelOut.QueueDeclare(c.qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error queue connection %s: %w", c.qName, err)
	}

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for result := range input {
			// time.Sleep(time.Second * 10)
			err := c.channelOut.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				MessageId:   "111",
				Timestamp:   time.Now(),
				Body:        []byte(result),
			})
			if err != nil {
				// send error type of result
			}
		}
	}()

	return done, nil
}

func (c *RMQOutputStream) Close() error {
	return c.channelOut.Close()
}
