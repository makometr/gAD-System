package rmq

import (
	"fmt"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/model"
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

func (c *RMQOutputStream) Result(input <-chan model.ResultWithID) (<-chan struct{}, error) {
	q, err := c.channelOut.QueueDeclare(c.qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error queue connection %s: %w", c.qName, err)
	}

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for msg := range input {
			var dataToSend []byte
			if dataToSend, err = msg.ToProto(); err != nil {
				fmt.Println("UNEXPECTED error while convert msg po proto!!!", msg)
			}

			// time.Sleep(time.Second * 10)
			err := c.channelOut.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				MessageId:   msg.ID,
				Timestamp:   time.Now(),
				Body:        dataToSend,
			})
			if err != nil {
				fmt.Println("UNEXPECTED error while send to out queue!!!", msg)
			}
		}
	}()

	return done, nil
}

func (c *RMQOutputStream) Close() error {
	return c.channelOut.Close()
}
