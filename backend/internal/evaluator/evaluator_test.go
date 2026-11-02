package evaluator

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	for _, tc := range []struct {
		expression string
		want       float64
	}{
		{"5+3", 8}, {"5-3", 2}, {"5*3", 15}, {"10/4", 2.5},
		{"5 + 3 * 2", 11}, {"10 - 6 / 2", 7}, {"8/2*3", 12}, {"10-3-2", 5},
		{"-3 * 2", -6}, {"2 * -3", -6}, {"2--3", 5}, {"--3", 3}, {"+3", 3},
		{"0", 0}, {"0*-5", 0}, {".5 + 1.", 1.5}, {"0.1+0.2", 0.3},
		{" \t 5\n+\r 3\u00a0", 8}, {"1e-7*2", 2e-7}, {"1E+2/4", 25},
		{strings.Repeat("-", 1023) + "1", -1}, {"1e-300*1e-300", 0},
	} {
		t.Run(tc.expression, func(t *testing.T) {
			got, err := Evaluate(tc.expression)
			if err != nil || math.Abs(got-tc.want) > 1e-12*math.Max(1, math.Abs(tc.want)) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	for _, tc := range []struct{ expression, code string }{
		{"", "invalid_expression"}, {" \t", "invalid_expression"}, {"5 + * 2", "invalid_expression"},
		{"5/0", "division_by_zero"}, {"1/-0", "division_by_zero"}, {"0/0", "division_by_zero"},
		{"1+", "invalid_expression"}, {"*2", "invalid_expression"}, {".", "invalid_expression"},
		{"2*", "invalid_expression"}, {"2/", "invalid_expression"}, {"2*/3", "invalid_expression"},
		{"1..2", "invalid_expression"}, {"1 2", "invalid_expression"}, {"(2+3)", "invalid_expression"},
		{"NaN", "invalid_expression"}, {"Infinity", "invalid_expression"}, {"2^3", "invalid_expression"},
		{"1e", "invalid_expression"}, {"1e+", "invalid_expression"}, {"0x10", "invalid_expression"},
		{"2%", "invalid_expression"}, {"1e999", "numeric_overflow"}, {strings.Repeat("9", 400), "numeric_overflow"},
		{"1e308*10", "numeric_overflow"}, {"1e308+1e308", "numeric_overflow"},
		{"-1e308-1e308", "numeric_overflow"}, {"1e308/1e-308", "numeric_overflow"},
		{strings.Repeat("1", 1025), "expression_too_large"},
	} {
		t.Run(tc.expression, func(t *testing.T) {
			_, err := Evaluate(tc.expression)
			var calculationError *Error
			if !errors.As(err, &calculationError) || calculationError.Code != tc.code || err.Error() == "" {
				t.Fatalf("got %v; want %s", err, tc.code)
			}
		})
	}
}

func FuzzEvaluate(f *testing.F) {
	for _, s := range []string{"1+2*3", "5/0", "--.5", "1e308*2", "", "\xff"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		value, err := Evaluate(s)
		if err == nil && (math.IsNaN(value) || math.IsInf(value, 0)) {
			t.Fatal("non-finite success")
		}
	})
}
