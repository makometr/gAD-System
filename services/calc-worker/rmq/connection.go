package rmq

import (
	"fmt"
	schema "gAD-System/internal/proto/expression/event"
	"gAD-System/services/calc-worker/config"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

// Connection represents connection vars to RMQ
type Connection struct {
	connection *amqp.Connection
	channelIn  *amqp.Channel
	channelOut *amqp.Channel

	producerResult *ResultProcucer

	consumerPlus  *OperationConsumer
	consumerMinus *OperationConsumer
	consumerMulti *OperationConsumer
	consumerMod   *OperationConsumer
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
		consumerPlus:  newOperationConsumer(rmqChanIn, cfg.RMQConfig.QNamePLus),
		consumerMinus: newOperationConsumer(rmqChanIn, cfg.RMQConfig.QNameMinus),
		consumerMulti: newOperationConsumer(rmqChanIn, cfg.RMQConfig.QNameMulti),
		consumerMod:   newOperationConsumer(rmqChanIn, cfg.RMQConfig.QNameMod),
	}, nil
}

type OperationConsumer struct {
	out chan<- string // expr to calculate type of channel
}

func newOperationConsumer(ch *amqp.Channel, qName string) *OperationConsumer {
	out := make(chan string) // expr result  type of channel
	exprs, err := ch.Consume(qName, "",
		true, false, false, false, nil,
	)
	if err != nil {
		// log fmt.Errorf("error consuming open: %w", err)
		return nil
	}

	go func() {
		for msg := range exprs {
			// раскодируем сообщение
			expr, _ := protoToMsg(msg.Body)

			// если ошибка, то что???
			// if err != nil {
			// return nil, fmt.Errorf("error proto to msg with id = %s: %w", msg.MessageId, err)
			// continue
			// }

			// отправляем раскодированное сообщение одному из воркеров через
			out <- string(expr)
		}
	}()

	return &OperationConsumer{out: out}
}

type ResultProcucer struct {
	in <-chan string // result of expr type of channel
}

func newResultProducer(ch *amqp.Channel, qName string) (*ResultProcucer, error) {
	in := make(chan string) // expr result type of channel

	go func() {

		for result := range in {
			fmt.Println(result)
			err := ch.Publish("", qName, false, false, amqp.Publishing{
				// ContentType: msg.ContentType,
				// MessageId:   msg.MessageId,
				// Timestamp:   time.Now(),
				// Body:        body,
			})
			if err != nil {
				// ????????????????????????
				// fmt.Errorf("error msg publishing with id = %s to queue = %s: %w", msg.MessageId, c.qNameOut, err)
			}
		}
	}()

	return &ResultProcucer{in: in}, nil
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
