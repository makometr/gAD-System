package rmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gAD-System/services/calc-controller/config"
	"gAD-System/services/calc-controller/model"

	"time"

	"github.com/streadway/amqp"
)

var ErrProtobuffSerialize = errors.New("error converting message to proto bytes")
var ErrAMQSend = errors.New("error seding msg to rmq")

type Producer interface {
	SendExpresion(ctx context.Context, expr model.ExprToCalc) error
	Close() error
}

type rmqProducer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	qNames map[model.Operation]string
}

func NewProducer(cfg *config.Config, connection *amqp.Connection) (Producer, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	var qNames = map[model.Operation]string{
		model.Plus:  cfg.RMQConfig.QNamePLus,
		model.Minus: cfg.RMQConfig.QNameMinus,
		model.Multi: cfg.RMQConfig.QNameMulti,
		model.Div:   cfg.RMQConfig.QNameDiv,
		model.Mod:   cfg.RMQConfig.QNameMod,
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

func (p *rmqProducer) SendExpresion(ctx context.Context, expr model.ExprToCalc) error {
	// serialize, err := MsgToProtoBytes(expr.Expr)
	serialize, err := json.Marshal(expr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProtobuffSerialize, err)
	}

	msg := amqp.Publishing{
		ContentType: "text/plain",
		MessageId:   string(expr.ID),
		Timestamp:   time.Now(),
		Body:        serialize,
	}

	err = p.ch.Publish("", p.qNames[expr.Oper], false, false, msg)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAMQSend, err)
	}

	fmt.Println("Sended to", p.qNames[expr.Oper])
	return nil
}

func (p *rmqProducer) Close() error {
	return p.ch.Close()
}
