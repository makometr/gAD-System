package rmq

import (
	"context"
	"fmt"
	"gAD-System/services/calc-controller/model"
	"strconv"
	"sync"
)

type ExprCalculator interface {
	CalculateExpression(model.Expression) (model.Result, error)
}

type RemoteCalculator struct {
	producer Producer
	consumer Consumer
	IDGEN    IDgenerator
}

func NewRemoteCalculator(prod Producer, cons Consumer) *RemoteCalculator {
	return &RemoteCalculator{producer: prod, consumer: cons, IDGEN: IDgenerator{Mutex: sync.Mutex{}}}
}

func (rc *RemoteCalculator) CalculateExpression(expr model.Expression) (model.Result, error) {
	ID := rc.IDGEN.GenereateID() // TODO normalno generate ID
	ctx := context.Background()

	recieveResult := make(chan model.ResultFromCalc)
	rc.consumer.Consume(ctx, recieveResult, ID)
	err := rc.producer.SendExpresion(ctx, model.ExprToCalc{Expression: expr, ID: ID})
	if err != nil {
		return model.Result{}, err
	}

	select {
	case result := <-recieveResult:
		fmt.Println("Received in calc expression:", result)
		return model.Result{result.Result.Result}, nil
	case <-ctx.Done():
		return model.Result{}, fmt.Errorf("calc error")

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
