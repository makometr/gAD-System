package worker

import (
	"context"
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

type WorkerConfig struct {
	input     <-chan ExprToCalc
	output    chan<- ResultExpr
	operation func(int64, int64) int64
	count     int
	delayGen  func() time.Duration
}

func StartWorker(ctx context.Context, cfg WorkerConfig) {
	for i := 0; i < cfg.count; i++ {
		go func() {
			for in := range cfg.input {
				time.Sleep(cfg.delayGen())
				cfg.output <- ResultExpr{result: cfg.operation(in.lhs, in.rhs), ID: in.ID}
			}
		}()
	}
}
