package model

type MsgID string

type Operation int

// TODO change this to map or to enum from intertnal/event/operation
const (
	Plus Operation = iota
	Minus
	Multi
	Div
	Mod
)

type Expression struct {
	Lhs  int64
	Rhs  int64
	Oper Operation
}

type ExprToCalc struct {
	Expression
	ID MsgID
}

type Result struct {
	Result int64
}

type ResultFromCalc struct {
	Result
	ID MsgID
}
