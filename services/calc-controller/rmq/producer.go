package rmq

import (
	"context"
	"errors"
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	"gAD-System/services/calc-controller/config"
	"gAD-System/services/calc-controller/model"

	"time"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

var ErrProtobuffSerialize = errors.New("error converting message to proto bytes")
var ErrAMQSend = errors.New("error seding msg to rmq")

type Producer interface {
	SendExpresion(ctx context.Context, expr pr_expr.Event, ID model.MsgID) error
	Close() error
}

type rmqProducer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	qNames map[pr_expr.Operation]string
}

func NewProducer(cfg *config.Config, connection *amqp.Connection) (Producer, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	var qNames = map[pr_expr.Operation]string{
		pr_expr.Operation_PLUS:  cfg.RMQConfig.QNamePLus,
		pr_expr.Operation_MINUS: cfg.RMQConfig.QNameMinus,
		pr_expr.Operation_MULTI: cfg.RMQConfig.QNameMulti,
		pr_expr.Operation_DIV:   cfg.RMQConfig.QNameDiv,
		pr_expr.Operation_MOD:   cfg.RMQConfig.QNameMod,
	}

	for _, qName := range qNames {
		_, err := channel.QueueDeclare(qName, true, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("error queue connection %s: %w", qName, err)
		}
	}

	producer := &rmqProducer{
		conn:   connection,
		ch:     channel,
		qNames: qNames,
	}

	return producer, nil
}

func (p *rmqProducer) SendExpresion(ctx context.Context, expr pr_expr.Event, ID model.MsgID) error {
	serialize, err := proto.Marshal(&expr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProtobuffSerialize, err)
	}

	msg := amqp.Publishing{
		ContentType: "text/plain",
		MessageId:   string(ID),
		Timestamp:   time.Now(),
		Body:        serialize,
	}

	err = p.ch.Publish("", p.qNames[expr.Operation], false, false, msg)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAMQSend, err)
	}

	fmt.Println("Sended to", p.qNames[expr.Operation])
	return nil
}

func (p *rmqProducer) Close() error {
	return p.ch.Close()
}
