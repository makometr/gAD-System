package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSimpleExpression(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErrType error
	}{
		{
			name: "Standart correct",
			args: args{expression: "150+50"},
			want: "200",
		},
		{
			name:        "Check error type-1",
			args:        args{expression: "150+++50"},
			wantErrType: ErrParseExpression,
		},
		{
			name:        "Check error type-2",
			args:        args{expression: "15+0*50"},
			wantErrType: ErrParseExpression,
		},
		{
			name:        "Check error type-3",
			args:        args{expression: "100/0"},
			wantErrType: ErrCalculationExpression,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateSimpleExpression(tt.args.expression)
			assert.ErrorIs(t, err, tt.wantErrType)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_parseResult_Calculate(t *testing.T) {
	tests := []struct {
		name    string
		pr      parseResult
		want    string
		wantErr bool
	}{
		{
			name: "plus test",
			pr:   parseResult{lhs: 10, rhs: 20, operation: plus},
			want: "30",
		},
		{
			name: "minus test",
			pr:   parseResult{lhs: 10, rhs: 20, operation: minus},
			want: "-10",
		},
		{
			name: "multiply test",
			pr:   parseResult{lhs: 25, rhs: 5, operation: multi},
			want: "125",
		},
		{
			name: "divide test",
			pr:   parseResult{lhs: 25, rhs: 5, operation: divide},
			want: "5",
		},
		{
			name: "divide with remainder test",
			pr:   parseResult{lhs: 100, rhs: 30, operation: divide},
			want: "3",
		},
		{
			name: "int64 test",
			pr:   parseResult{lhs: -9223372036854775808, rhs: 9223372036854775807, operation: plus},
			want: "-1",
		},

		{
			name:    "incorrect dividing by zero",
			pr:      parseResult{lhs: 9876, rhs: 0, operation: divide},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pr.calculate()
			require.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_parseBinaryExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    parseResult
		wantErr bool
	}{
		// correct
		{
			name: "correct plus",
			args: args{expr: "100+100"},
			want: parseResult{lhs: 100, rhs: 100, operation: plus},
		},
		{
			name: "correct minus",
			args: args{expr: "100-1"},
			want: parseResult{lhs: 100, rhs: 1, operation: minus},
		},
		{
			name: "correct multi",
			args: args{expr: "2*10000"},
			want: parseResult{lhs: 2, rhs: 10000, operation: multi},
		},
		{
			name: "correct divide",
			args: args{expr: "0/20"},
			want: parseResult{lhs: 0, rhs: 20, operation: divide},
		},
		{
			name: "correct divide by zero",
			args: args{expr: "20/0"},
			want: parseResult{lhs: 20, rhs: 0, operation: divide},
		},

		// incorrect
		{
			name:    "incorrect lhs",
			args:    args{expr: "10f0+200"},
			wantErr: true,
		},
		{
			name:    "incorrect rhs",
			args:    args{expr: "100+20ggg0"},
			wantErr: true,
		},
		{
			name:    "incorrect no operation",
			args:    args{expr: "100g200"},
			wantErr: true,
		},
		{
			name:    "incorrect many operations",
			args:    args{expr: "100+*+200"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBinaryExpression(tt.args.expr)
			require.Equal(t, err != nil, tt.wantErr)
			assert.Exactly(t, got, tt.want)
		})
	}
}
