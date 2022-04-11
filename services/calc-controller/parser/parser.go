package parser

import (
	"errors"
	"fmt"
	pr_expr "gAD-System/internal/proto/expression/event"
	"strconv"
	"strings"
)

var (
	// ErrParseExpression occurres when expression is not like "NumberOperationNumber"
	ErrParseExpression = errors.New("error in parsing expression")

	// ErrCalculationExpression occurres with incorrect arithmetic operations, divide by zero for example
	ErrCalculationExpression = errors.New("error in calculations")
)

var operations = map[rune]pr_expr.Operation{
	'+': pr_expr.Operation_PLUS,
	'-': pr_expr.Operation_MINUS,
	'*': pr_expr.Operation_MULTI,
	'/': pr_expr.Operation_DIV,
	'%': pr_expr.Operation_MOD,
}

func ParseBinaryExpression(expr string) (*pr_expr.Event, error) {
	operIndex := strings.IndexFunc(expr, func(r rune) bool {
		for sign := range operations {
			if r == sign {
				return true
			}
		}
		return false
	})
	if operIndex == -1 {
		return nil, fmt.Errorf("no operation found")
	}

	lhsNum, err := strconv.ParseInt(expr[:operIndex], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("left number invalid: %s", expr[:operIndex])
	}
	rhsNum, err := strconv.ParseInt(expr[operIndex+1:], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("right number invalid: %s", expr[operIndex+1:])
	}

	return &pr_expr.Event{Lhs: lhsNum, Rhs: rhsNum, Operation: operations[rune(expr[operIndex])]}, nil
}
