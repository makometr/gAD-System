package rmq

import (
	"fmt"
	"gAD-System/services/calc-worker/config"

	"github.com/streadway/amqp"
)

type RMQInputStream struct {
	channelIn *amqp.Channel
	qNames    map[Operation]string
}

func NewRMQInputStream(cfg *config.Config, rmq *Connection) (InputExprStream, error) {
	rmqChanIn, err := rmq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel error: %w", err)
	}
	var queueNames = map[Operation]string{
		Plus:  cfg.RMQConfig.QNamePLus,
		Minus: cfg.RMQConfig.QNameMinus,
		Multi: cfg.RMQConfig.QNameMulti,
		Div:   cfg.RMQConfig.QNameDiv,
		Mod:   cfg.RMQConfig.QNameMod,
	}
	return &RMQInputStream{channelIn: rmqChanIn, qNames: queueNames}, nil
}

func (c *RMQInputStream) Plus() (<-chan string, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[Plus])
}

func (c *RMQInputStream) Minus() (<-chan string, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[Minus])
}

func (c *RMQInputStream) Milti() (<-chan string, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[Multi])
}

func (c *RMQInputStream) Div() (<-chan string, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[Div])
}

func (c *RMQInputStream) Mod() (<-chan string, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[Mod])
}

func (c *RMQInputStream) Close() error {
	return c.channelIn.Close()
}

func newExpressionConsumer(ch *amqp.Channel, qName string) (chan string, error) {
	q, err := ch.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error queue connection %s: %w", qName, err)
	}

	toCalc := make(chan string) // expr result type of channel
	exprs, err := ch.Consume(q.Name, "",
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error consuming open: %w", err)
	}

	go func() {
		for msg := range exprs {
			// раскодируем сообщение
			// expr, _ := protoToMsg(msg.Body)

			// если ошибка, то что???
			// if err != nil {
			// return nil, fmt.Errorf("error proto to msg with id = %s: %w", msg.MessageId, err)
			// continue
			// }

			// отправляем раскодированное сообщение одному из воркеров через канал
			toCalc <- string(msg.Body)
		}
		close(toCalc)
	}()

	return toCalc, nil
}
