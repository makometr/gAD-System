package worker

import (
	"context"
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	pr_result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-worker/config"
	"gAD-System/services/calc-worker/model"
	"gAD-System/services/calc-worker/rmq"
	"sync"
	"time"
)

// type ExprToCalc struct {
// 	lhs int64
// 	rhs int64
// 	ID  string
// }

// type ResultExpr struct {
// 	result int64
// 	ID     string
// }

type calcFunction func(int64, int64) int64

var operations = map[pr_expr.Operation]calcFunction{
	pr_expr.Operation_PLUS:  func(lhs, rhs int64) int64 { return lhs + rhs },
	pr_expr.Operation_MINUS: func(lhs, rhs int64) int64 { return lhs - rhs },
	pr_expr.Operation_MULTI: func(lhs, rhs int64) int64 { return lhs * rhs },
	pr_expr.Operation_DIV:   func(lhs, rhs int64) int64 { return lhs / rhs },
	pr_expr.Operation_MOD:   func(lhs, rhs int64) int64 { return lhs % rhs },
}

type WorkerConfig struct {
	Input     <-chan model.ExprWithID
	Output    chan<- model.ResultWithID
	Operation pr_expr.Operation
	DelayGen  func() time.Duration
}

func StartWorkers(ctx context.Context, cfg *config.Config, input rmq.InputExprStream, output rmq.OutputExprStream) (<-chan struct{}, error) {
	workersGroup := sync.WaitGroup{}
	done := make(chan struct{})

	chResult := make(chan model.ResultWithID, 10) // TODO expr result type of channel
	sendDone, err := output.Result(chResult)
	if err != nil {
		return nil, fmt.Errorf("cant init out stream: %w", err)
	}

	chPlus, err := input.Plus()
	if err != nil {
		return nil, fmt.Errorf("cant init in-plus stream: %w", err)
	}
	plusCfg := WorkerConfig{Input: chPlus, Output: chResult, Operation: pr_expr.Operation_PLUS,
		DelayGen: nil,
		// DelayGen: func() time.Duration { return time.Second * 10 },
	}
	startOperationWorker(plusCfg, cfg.WConfig.CountPlus, &workersGroup)

	chMinus, err := input.Minus()
	if err != nil {
		return nil, fmt.Errorf("cant init in-minus stream: %w", err)
	}
	minusCfg := WorkerConfig{Input: chMinus, Output: chResult, Operation: pr_expr.Operation_MINUS,
		DelayGen: nil,
		// DelayGen: func() time.Duration { return time.Second * 10 },
	}
	startOperationWorker(minusCfg, cfg.WConfig.CountPlus, &workersGroup)

	chMulti, err := input.Milti()
	if err != nil {
		return nil, fmt.Errorf("cant init in-minus stream: %w", err)
	}
	multiCfg := WorkerConfig{Input: chMulti, Output: chResult, Operation: pr_expr.Operation_MULTI,
		DelayGen: nil,
		// DelayGen: func() time.Duration { return time.Second * 10 },
	}
	startOperationWorker(multiCfg, cfg.WConfig.CountPlus, &workersGroup)

	chDiv, err := input.Div()
	if err != nil {
		return nil, fmt.Errorf("cant init in-minus stream: %w", err)
	}
	divCfg := WorkerConfig{Input: chDiv, Output: chResult, Operation: pr_expr.Operation_DIV,
		DelayGen: nil,
		// DelayGen: func() time.Duration { return time.Second * 10 },
	}
	startOperationWorker(divCfg, cfg.WConfig.CountPlus, &workersGroup)

	chMod, err := input.Mod()
	if err != nil {
		return nil, fmt.Errorf("cant init in-minus stream: %w", err)
	}
	modCfg := WorkerConfig{Input: chMod, Output: chResult, Operation: pr_expr.Operation_MOD,
		DelayGen: nil,
		// DelayGen: func() time.Duration { return time.Second * 10 },
	}
	startOperationWorker(modCfg, cfg.WConfig.CountPlus, &workersGroup)

	go func() {
		workersGroup.Wait() // Воркеры завершили работу и передали ответы на отправку
		close(chResult)     // Закрывем канал отправки результатов, сигнализируя, что новых не будет
		<-sendDone          // ждём уведомления, что все переданные на отправку сообщения отправлены
		done <- struct{}{}
	}()

	return done, nil
}

type Result struct {
	Result int64
}

type ResultFromCalc struct {
	Result
	ID string
}

func startOperationWorker(cfg WorkerConfig, count int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for in := range cfg.Input {
				if cfg.DelayGen != nil {
					time.Sleep(cfg.DelayGen())
				}
				ans := calculate(in.Expr.Lhs, in.Expr.Rhs, cfg.Operation)
				cfg.Output <- model.ResultWithID{Result: pr_result.Event{Result: &pr_result.Event_Product{Product: ans}}, ID: in.ID}
			}
		}()
	}
}

func calculate(lhs, rhs int64, operation pr_expr.Operation) int64 {
	return operations[operation](lhs, rhs)
}
