package rmq

import (
	"fmt"
	"gAD-System/services/calc-worker/config"

	"github.com/streadway/amqp"
)

type Operation int

// TODO change this to map or to enum from intertnal/event/operation
const (
	Plus Operation = iota
	Minus
	Multi
	Div
	Mod
)

// Connection represents connection vars to RMQ
type Connection struct {
	conn *amqp.Connection
}

// InitConnection inits only raw connection with RMQ
func InitConnection(cfg *config.Config) (*Connection, error) {
	rmqConn, err := amqp.Dial(fmt.Sprintf("amqp://%s", cfg.RMQConfig.Server))
	if err != nil {
		return nil, fmt.Errorf("init connection error: %w", err)
	}

	return &Connection{conn: rmqConn}, nil
}

// Close closes connection to rmq, to graceful shutdown
func (c *Connection) Close() error {
	return c.conn.Close()
}

// func msgToProtoBytes(message string) ([]byte, error) {
// 	event := schema.Event{Expression: message}
// 	out, err := proto.Marshal(&event)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }

// func protoToMsg(data []byte) (string, error) {
// 	event := schema.Event{}
// 	if err := proto.Unmarshal(data, &event); err != nil {
// 		return "", err
// 	}
// 	return event.Expression, nil
// }
