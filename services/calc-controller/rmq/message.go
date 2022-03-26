package rmq

import "time"

type Message struct {
	ContentType string
	Timestamp   time.Time
	MessageID   string
	Body        []byte
}
