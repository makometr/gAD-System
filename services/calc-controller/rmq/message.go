package rmq

import (
	"time"

	schema "gAD-System/internal/proto/expression/event"

	"google.golang.org/protobuf/proto"
)

type MsgID string

type Message struct {
	ContentType string
	Timestamp   time.Time
	MessageID   MsgID
	Body        []byte
}

// ExpressionWithID linked raw expression with id
type ExpressionWithID struct {
	Expr string
	Id   MsgID
}

func (e ExpressionWithID) ToProto() ([]byte, error) {
	// event := &schema.Event{Expression: e.expr}
	// out, err := proto.Marshal(event)
	// if err != nil {
	// 	return nil, err
	// }
	// return out, nil

	return nil, nil
}

func MsgToProtoBytes(message string) ([]byte, error) {
	event := schema.Event{Expression: message}
	out, err := proto.Marshal(&event)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ProtoToMsg(data []byte) (string, error) {
	event := schema.Event{}
	if err := proto.Unmarshal(data, &event); err != nil {
		return "", err
	}
	return event.Expression, nil
}
