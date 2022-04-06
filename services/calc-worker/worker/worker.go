package worker

import (
	"context"
	"fmt"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/rmq"
	"sync"
	"time"
)

type ExprToCalc struct {
	lhs int64
	rhs int64
	ID  string
}

type ResultExpr struct {
	result int64
	ID     string
}

type calcFunction func(int64, int64) int64

var operations = map[rmq.Operation]calcFunction{
	rmq.Plus:  func(lhs, rhs int64) int64 { return lhs + rhs },
	rmq.Minus: func(lhs, rhs int64) int64 { return lhs - rhs },
	rmq.Multi: func(lhs, rhs int64) int64 { return lhs * rhs },
	rmq.Div:   func(lhs, rhs int64) int64 { return lhs / rhs },
	rmq.Mod:   func(lhs, rhs int64) int64 { return lhs % rhs },
}

type WorkerConfig struct {
	Input     <-chan string // ExprToCalc
	Output    chan<- string // ResultExpr
	Operation rmq.Operation
	DelayGen  func() time.Duration
}

func StartWorkers(ctx context.Context, cfg *config.Config, input rmq.InputExprStream, output rmq.OutputExprStream) (<-chan struct{}, error) {
	workersGroup := sync.WaitGroup{}
	done := make(chan struct{})

	chResult := make(chan string, 10) // TODO expr result type of channel
	sendDone, err := output.Result(chResult)
	if err != nil {
		return nil, fmt.Errorf("cant init out stream: %w", err)
	}

	chPlus, err := input.Plus()
	if err != nil {
		return nil, fmt.Errorf("cant init in-plus stream: %w", err)
	}
	for i := 0; i < cfg.WConfig.CountPlus; i++ {
		workersGroup.Add(1)
		go func() {
			startWorker(ctx, WorkerConfig{
				Input:     chPlus,
				Output:    chResult,
				Operation: rmq.Plus,
				DelayGen:  func() time.Duration { return time.Second * 10 },
			})
			workersGroup.Done()
		}()
	}

	chMinus, err := input.Minus()
	if err != nil {
		return nil, fmt.Errorf("cant init in-minus stream: %w", err)
	}
	for i := 0; i < cfg.WConfig.CountMinus; i++ {
		workersGroup.Add(1)
		go func() {
			startWorker(ctx, WorkerConfig{
				Input:     chMinus,
				Output:    chResult,
				Operation: rmq.Minus,
				DelayGen:  nil,
			})
			workersGroup.Done()
		}()
	}

	go func() {
		workersGroup.Wait() // Воркеры завершили работу и передали ответы на отправку
		close(chResult)     // Закрывем канал отправки результатов, сигнализируя, что новых не будет
		<-sendDone          // ждём уведомления, что все переданные на отправку сообщения отправлены
		done <- struct{}{}
	}()

	return done, nil
}

func startWorker(ctx context.Context, cfg WorkerConfig) {
	for in := range cfg.Input {
		if cfg.DelayGen != nil {
			time.Sleep(cfg.DelayGen())
		}
		// cfg.Output <- ResultExpr{result: operations[cfg.Operation](in.lhs, in.rhs), ID: in.ID}
		cfg.Output <- in + "+"
	}
}
