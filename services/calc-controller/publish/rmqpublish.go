package publish

import "github.com/streadway/amqp"

type RMQPublisher struct {
	conn *amqp.Connection
}

func NewRMQPublisher() Publisher {
	return &RMQPublisher{}
}

func (p *RMQPublisher) Publish() {

}
