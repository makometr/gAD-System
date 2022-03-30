package rmq

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/streadway/amqp"
)

const fakeAddress = "amqp://localhost:5672/"

var (
	fakeServer     *server.AMQPServer
	mockConnection *amqp.Connection
)

func TestRmqPublisher_Publish(t *testing.T) {
	// configureEnv()
	// publisher := NewPublisher(mockConnection, "")
	// input := make(chan Message)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		input <- Message{
	// 			ContentType: "text/plain",
	// 			Timestamp:   time.Now(),
	// 			MessageID:   "MsgID",
	// 			Body:        nil,
	// 		}
	// 	}
	// 	close(input)
	// }()
	// fmt.Println("start publishing")
	// err := publisher.Publish(context.Background(), input)
	// assert.NoError(t, err)
}

func TestNewPublisher(t *testing.T) {

}

func configureEnv() {
	fakeServer = server.NewServer(fakeAddress)
	fakeServer.Start()

	mockConnection, _ = amqp.Dial(fakeAddress)
}
