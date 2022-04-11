package rmq

import (
	"context"
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	pr_result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-controller/model"
	"strconv"
	"sync"
)

type ExprCalculator interface {
	CalculateExpression(pr_expr.Event) (pr_result.Event, error)
}

type RemoteCalculator struct {
	producer Producer
	consumer Consumer
	IDGEN    IDgenerator
}

func NewRemoteCalculator(prod Producer, cons Consumer) *RemoteCalculator {
	return &RemoteCalculator{producer: prod, consumer: cons, IDGEN: IDgenerator{Mutex: sync.Mutex{}}}
}

func (rc *RemoteCalculator) CalculateExpression(expr pr_expr.Event) (pr_result.Event, error) {
	ID := rc.IDGEN.GenereateID() // TODO normalno generate ID
	ctx := context.Background()

	recieveResult := make(chan pr_result.Event)
	rc.consumer.Consume(ctx, recieveResult, ID)
	err := rc.producer.SendExpresion(ctx, expr, ID)
	if err != nil {
		return pr_result.Event{}, err
	}

	select {
	case msg := <-recieveResult:
		fmt.Println("Received in calc expression:", msg)

		switch res := msg.Result.(type) {
		case *pr_result.Event_ErrorMsg:
			return pr_result.Event{}, fmt.Errorf(res.ErrorMsg)
		case *pr_result.Event_Product:
			return msg, nil
		default:
			return pr_result.Event{}, fmt.Errorf("enexpected value from protobuf convestion")
		}

		return msg, nil
	case <-ctx.Done():
		return pr_result.Event{}, fmt.Errorf("calc error")
	}
}

type IDgenerator struct {
	counter int
	sync.Mutex
}

func (idg *IDgenerator) GenereateID() model.MsgID {
	idg.Lock()
	defer idg.Unlock()
	defer func() {
		idg.counter++
	}()

	return model.MsgID(strconv.Itoa(idg.counter))
}
