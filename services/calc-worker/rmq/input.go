package rmq

import (
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/model"

	"github.com/streadway/amqp"
)

type RMQInputStream struct {
	channelIn *amqp.Channel
	qNames    map[pr_expr.Operation]string
}

func NewRMQInputStream(cfg *config.Config, rmq *Connection) (InputExprStream, error) {
	rmqChanIn, err := rmq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel error: %w", err)
	}
	var queueNames = map[pr_expr.Operation]string{
		pr_expr.Operation_PLUS:  cfg.RMQConfig.QNamePLus,
		pr_expr.Operation_MINUS: cfg.RMQConfig.QNameMinus,
		pr_expr.Operation_MULTI: cfg.RMQConfig.QNameMulti,
		pr_expr.Operation_DIV:   cfg.RMQConfig.QNameDiv,
		pr_expr.Operation_MOD:   cfg.RMQConfig.QNameMod,
	}
	return &RMQInputStream{channelIn: rmqChanIn, qNames: queueNames}, nil
}

func (c *RMQInputStream) Plus() (<-chan model.ExprWithID, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[pr_expr.Operation_PLUS])
}

func (c *RMQInputStream) Minus() (<-chan model.ExprWithID, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[pr_expr.Operation_MINUS])
}

func (c *RMQInputStream) Milti() (<-chan model.ExprWithID, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[pr_expr.Operation_MULTI])
}

func (c *RMQInputStream) Div() (<-chan model.ExprWithID, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[pr_expr.Operation_DIV])
}

func (c *RMQInputStream) Mod() (<-chan model.ExprWithID, error) {
	return newExpressionConsumer(c.channelIn, c.qNames[pr_expr.Operation_MOD])
}

func (c *RMQInputStream) Close() error {
	return c.channelIn.Close()
}

func newExpressionConsumer(ch *amqp.Channel, qName string) (chan model.ExprWithID, error) {
	q, err := ch.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error queue connection %s: %w", qName, err)
	}

	toCalc := make(chan model.ExprWithID) // expr result type of channel
	exprs, err := ch.Consume(q.Name, "",
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error consuming open: %w", err)
	}

	go func() {
		for msg := range exprs {
			expr, err := model.NewExprWithIDFromBytes(msg.Body, msg.MessageId)
			if err != nil {
				// ???????
			}

			toCalc <- *expr
		}
		close(toCalc)
	}()

	return toCalc, nil
}
