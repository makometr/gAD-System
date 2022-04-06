package rmq

// TODO change type oof chans to expr ans results

type InputExprStream interface {
	Plus() (<-chan string, error)
	Minus() (<-chan string, error)
	Milti() (<-chan string, error)
	Div() (<-chan string, error)
	Mod() (<-chan string, error)

	Close() error
}

type OutputExprStream interface {
	Result(<-chan string) (<-chan struct{}, error)

	Close() error
}
