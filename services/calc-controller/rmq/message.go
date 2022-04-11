package rmq

import (
	"time"

	"gAD-System/services/calc-controller/model"
)

// type MsgID string

type MessageFromRMQ struct {
	ContentType string
	Timestamp   time.Time
	MessageID   model.MsgID
	Body        []byte
}
