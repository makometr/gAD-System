package rmq

import "gAD-System/services/calc-worker/model"

// TODO change type oof chans to expr ans results

type InputExprStream interface {
	Plus() (<-chan model.ExprWithID, error)
	Minus() (<-chan model.ExprWithID, error)
	Milti() (<-chan model.ExprWithID, error)
	Div() (<-chan model.ExprWithID, error)
	Mod() (<-chan model.ExprWithID, error)

	Close() error
}

type OutputExprStream interface {
	Result(<-chan model.ResultWithID) (<-chan struct{}, error)

	Close() error
}
