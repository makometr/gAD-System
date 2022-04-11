package worker

import (
	pr_expr "gAD-System/internal/proto/expression/event"
	"testing"
)

func Test_calculate(t *testing.T) {
	type args struct {
		lhs       int64
		rhs       int64
		operation pr_expr.Operation
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "standart plus",
			args: args{lhs: 100, rhs: 100, operation: pr_expr.Operation_PLUS},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculate(tt.args.lhs, tt.args.rhs, tt.args.operation); got != tt.want {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
