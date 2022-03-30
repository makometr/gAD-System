package rmq

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type ExprCalculator interface {
	CalculateExpression(string) (string, error)
}

type RemoteCalculator struct {
	producer Producer
	consumer Consumer
	IDGEN    IDgenerator
}

func NewRemoteCalculator(prod Producer, cons Consumer) *RemoteCalculator {
	return &RemoteCalculator{producer: prod, consumer: cons, IDGEN: IDgenerator{Mutex: sync.Mutex{}}}
}

func (rc *RemoteCalculator) CalculateExpression(expr string) (string, error) {
	ID := rc.IDGEN.GenereateID() // TODO normalno generate ID
	ctx := context.Background()

	recieveResult := make(chan ExpressionWithID)
	rc.consumer.Consume(ctx, recieveResult, MsgID(ID))
	err := rc.producer.SendExpresion(ctx, ExpressionWithID{Expr: expr, Id: MsgID(ID)})
	if err != nil {
		return "", err
	}

	select {
	case result := <-recieveResult:
		fmt.Println("Received in calc expression:", result)
		return result.Expr, nil
	case <-ctx.Done():
		return "", fmt.Errorf("calc error")

	}
}

type IDgenerator struct {
	counter int
	sync.Mutex
}

func (idg *IDgenerator) GenereateID() string {
	idg.Lock()
	defer idg.Unlock()
	defer func() {
		idg.counter++
	}()

	return strconv.Itoa(idg.counter)
}
