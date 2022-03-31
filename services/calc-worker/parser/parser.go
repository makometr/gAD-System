package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrParseExpression occurres when expression is not like "NumberOperationNumber"
	ErrParseExpression = errors.New("error in parsing expression")

	// ErrCalculationExpression occurres with incorrect arithmetic operations, divide by zero for example
	ErrCalculationExpression = errors.New("error in calculations")
)

// CalculateSimpleExpression returns result of arithmetic expression
func CalculateSimpleExpression(expression string) (string, error) {
	parsed, err := parseBinaryExpression(expression)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParseExpression, err)
	}

	result, err := parsed.calculate()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCalculationExpression, err)
	}

	return result, nil
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
		for sign := range operations {
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

func (pr parseResult) calculate() (string, error) {
	var result operandType

	switch pr.operation {
	case plus:
		result = pr.lhs + pr.rhs
	case minus:
		result = pr.lhs - pr.rhs
	case multi:
		result = pr.lhs * pr.rhs
	case divide:
		if pr.rhs == 0 {
			return "", fmt.Errorf("dividing by zero")
		}
		result = pr.lhs / pr.rhs
	}

	return strconv.FormatInt(int64(result), 10), nil
}
