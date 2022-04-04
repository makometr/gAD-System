package rmq

import (
	"fmt"
	schema "gAD-System/internal/proto/expression/event"
	"gAD-System/services/calc-worker/config"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

// Connection represents connection vars to RMQ
type Connection struct {
	connection *amqp.Connection
	channelIn  *amqp.Channel
	channelOut *amqp.Channel
	qNameIn    string
	qNameOut   string
}

// InitConnection inits connections and in-out channels with RMQ with provided config
func InitConnection(cfg *config.Config) (*Connection, error) {
	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://%s", cfg.RMQConfig.Server))
	if err != nil {
		return nil, fmt.Errorf("init connection error: %w", err)
	}

	rmqChanIn, err := rmqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel error: %w", err)
	}

	rmqChanOut, err := rmqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel error: %w", err)
	}

	return &Connection{connection: rmqConn, channelIn: rmqChanIn, channelOut: rmqChanOut,
		qNameIn: cfg.RMQConfig.PubQueryName, qNameOut: cfg.RMQConfig.SubQueryName}, nil
}

// Close closes connection to rmq, to graceful shutdown
func (c *Connection) Close() error {
	if err := c.channelIn.Close(); err != nil {
		return err
	}
	if err := c.channelOut.Close(); err != nil {
		return err
	}
	if err := c.connection.Close(); err != nil {
		return err
	}
	return nil
}

// CalculateExpressions calculation worker for go-call style.
// Depends on open chans from RMQ, no need to close by hands. Worker exited when connection closed.
// TODO error handling
func (c *Connection) CalculateExpressions(errChan chan<- error, calculator func(string) (string, error)) {
	exprs, err := c.channelIn.Consume(c.qNameIn, "",
		true, false, false, false, nil,
	)
	if err != nil {
		errChan <- fmt.Errorf("error consuming open: %w", err)
		return
	}

	for msg := range exprs {
		expr, err := protoToMsg(msg.Body)
		if err != nil {
			errChan <- fmt.Errorf("error proto to msg with id = %s: %w", msg.MessageId, err)
			// continue
		}

		result, err := calculator(expr)
		if err != nil {
			errChan <- fmt.Errorf("error calculator msg with id = %s: %w", msg.MessageId, err)
			// continue
		}

		body, err := msgToProtoBytes(result)
		if err != nil {
			errChan <- fmt.Errorf("error msg with id = %s to proto: %w", msg.MessageId, err)
			// continue
		}

		err = c.channelOut.Publish("", c.qNameOut, false, false, amqp.Publishing{
			ContentType: msg.ContentType,
			MessageId:   msg.MessageId,
			Timestamp:   time.Now(),
			Body:        body,
		})
		if err != nil {
			errChan <- fmt.Errorf("error msg publishing with id = %s to queue = %s: %w", msg.MessageId, c.qNameOut, err)
			continue
		}
	}
}

func msgToProtoBytes(message string) ([]byte, error) {
	event := schema.Event{Expression: message}
	out, err := proto.Marshal(&event)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func protoToMsg(data []byte) (string, error) {
	event := schema.Event{}
	if err := proto.Unmarshal(data, &event); err != nil {
		return "", err
	}
	return event.Expression, nil
}
