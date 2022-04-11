package model

import (
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	pr_result "gAD-System/internal/proto/result/event"

	"google.golang.org/protobuf/proto"
)

// decalre interface ith methods

// ExprWithID incapsulates proto-events and IDs for RMQ
type ExprWithID struct {
	Expr pr_expr.Event
	ID   string
}

// NewExprWithIDFromBytes creates new instance using proto lib
func NewExprWithIDFromBytes(data []byte, id string) (*ExprWithID, error) {
	ewi := ExprWithID{ID: id}
	if err := proto.Unmarshal(data, &ewi.Expr); err != nil {
		return &ewi, fmt.Errorf("convert error: %w", err)
	}
	return &ewi, nil
}

// ResultWithID incapsulates proto-events and IDs for RMQ
type ResultWithID struct {
	Result pr_result.Event
	ID     string
}

// ToProto converts to proto bytes using proto lib
func (rwi *ResultWithID) ToProto() ([]byte, error) {
	return proto.Marshal(&rwi.Result)
}
