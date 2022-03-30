package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrParseExpression = errors.New("error in parsing expression")
)

func CalculateSimpleExpression(expression string) (string, error) {
	parsed, err := parseBinaryExpression(expression)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParseExpression, err)
	}

	return parsed.calculate(), nil
}

type operation int
type operandType int64

const (
	plus operation = iota
	minus
	multi
	divide
)

var operations = map[rune]operation{
	'+': plus,
	'-': minus,
	'*': multi,
	'/': divide,
}

type parseResult struct {
	lhs       operandType
	rhs       operandType
	operation operation
}

func parseBinaryExpression(expr string) (parseResult, error) {
	operIndex := strings.IndexFunc(expr, func(r rune) bool {
		for sign, _ := range operations {
			if r == sign {
				return true
			}
		}
		return false
	})
	if operIndex == -1 {
		return parseResult{}, fmt.Errorf("no operation found")
	}

	lhsNum, err := strconv.ParseInt(expr[:operIndex], 10, 64)
	if err != nil {
		return parseResult{}, fmt.Errorf("left number invalid: %s", expr[:operIndex])
	}
	rhsNum, err := strconv.ParseInt(expr[operIndex+1:], 10, 64)
	if err != nil {
		return parseResult{}, fmt.Errorf("right number invalid: %s", expr[operIndex+1:])
	}

	return parseResult{lhs: operandType(lhsNum), rhs: operandType(rhsNum), operation: operations[rune(expr[operIndex])]}, nil
}

func (pr parseResult) calculate() string {
	var result operandType

	switch pr.operation {
	case plus:
		result = pr.lhs + pr.rhs
	case minus:
		result = pr.lhs - pr.rhs
	case multi:
		result = pr.lhs * pr.rhs
	case divide:
		result = pr.lhs / pr.rhs
	}

	return strconv.FormatInt(int64(result), 10)
}
