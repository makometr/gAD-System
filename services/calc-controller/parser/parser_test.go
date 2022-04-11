package parser

import (
	pr_expr "gAD-System/internal/proto/expression/event"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseBinaryExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    *pr_expr.Event
		wantErr bool
	}{
		// correct
		{
			name: "correct plus",
			args: args{expr: "100+100"},
			want: &pr_expr.Event{Lhs: 100, Rhs: 100, Operation: pr_expr.Operation_PLUS},
		},
		{
			name: "correct minus",
			args: args{expr: "100-1"},
			want: &pr_expr.Event{Lhs: 100, Rhs: 1, Operation: pr_expr.Operation_MINUS},
		},
		{
			name: "correct multi",
			args: args{expr: "2*10000"},
			want: &pr_expr.Event{Lhs: 2, Rhs: 10000, Operation: pr_expr.Operation_MULTI},
		},
		{
			name: "correct divide",
			args: args{expr: "0/20"},
			want: &pr_expr.Event{Lhs: 0, Rhs: 20, Operation: pr_expr.Operation_DIV},
		},
		{
			name: "correct mod",
			args: args{expr: "100%30"},
			want: &pr_expr.Event{Lhs: 100, Rhs: 30, Operation: pr_expr.Operation_MOD},
		},
		// {
		// 	name: "correct divide by zero",
		// 	args: args{expr: "20/0"},
		// 	want: &pr_expr.Event{Lhs: 20, Rhs: 0, Operation: pr_expr.Operation_DIV},
		// },

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
			got, err := ParseBinaryExpression(tt.args.expr)
			require.Equal(t, err != nil, tt.wantErr)
			assert.Exactly(t, got, tt.want)
		})
	}
}
