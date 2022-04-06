package rmq

import (
	"time"

	schema "gAD-System/internal/proto/expression/event"
	"gAD-System/services/calc-controller/model"

	"google.golang.org/protobuf/proto"
)

// type MsgID string

type Message struct {
	ContentType string
	Timestamp   time.Time
	MessageID   model.MsgID
	Body        []byte
}

// ExpressionWithID linked raw expression with id
type ExpressionWithID struct {
	Expr string
	Id   model.MsgID
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
