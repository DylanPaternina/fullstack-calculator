package evaluator

import (
	"errors"
	"math"
	"strings"
	"testing"
)

// Never let an absolute tolerance erase tiny values or accept non-finite results.
// Exact zero also protects deliberately documented underflow-to-zero cases.
func closeEnough(got, want float64) bool {
	if math.IsNaN(got) || math.IsNaN(want) || math.IsInf(got, 0) || math.IsInf(want, 0) {
		return false
	}
	if want == 0 {
		return got == 0
	}
	return math.Abs((got-want)/want) <= 1e-12
}

func TestCloseEnough(t *testing.T) {
	for _, tc := range []struct {
		name      string
		got, want float64
		match     bool
	}{
		{"NaN result", math.NaN(), 11, false},
		{"positive infinity", math.Inf(1), 11, false},
		{"negative infinity", math.Inf(-1), 11, false},
		{"NaN expectation", 1, math.NaN(), false},
		{"infinite expectation", 1, math.Inf(1), false},
		{"incorrect near zero", 1e-13, 0, false},
		{"subnormal is not zero", math.SmallestNonzeroFloat64, 0, false},
		{"zero is not subnormal", 0, math.SmallestNonzeroFloat64, false},
		{"different subnormals", 2 * math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, false},
		{"same subnormal", math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{"tiny relative error", 1.01e-300, 1e-300, false},
		{"ordinary decimal rounding", 0.30000000000000004, 0.3, true},
		{"exact zero", 0, 0, true},
		{"signed zero", math.Copysign(0, -1), 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := closeEnough(tc.got, tc.want); got != tc.match {
				t.Fatalf("comparison returned %v; want %v", got, tc.match)
			}
		})
	}
}

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
		{"5e-324", math.SmallestNonzeroFloat64}, {"1e-300/2", 5e-301},
		{"2+3*4", 14}, {"(2+3)*4", 20}, {"(2+3)", 5},
		{"2^3", 8}, {"5^2", 25}, {"10^0", 1}, {"0^0", 1}, {"2^-3", 0.125},
		{"2^3^2", 512}, {"-2^2", -4}, {"(-2)^2", 4}, {"2+3^2", 11}, {"2*3^2", 18},
		{"2^-3^2", 1.0 / 512}, {"(-2)^-3", -0.125}, {"2^--3", 8},
		{"sqrt(81)", 9}, {"sqrt(16)", 4}, {"sqrt(0)", 0}, {"sqrt(2)", math.Sqrt(2)},
		{"sqrt(9)+2", 5}, {"sqrt(25)*2", 10}, {"2^3+sqrt(16)", 12},
		{"sqrt(sqrt(256))", 4}, {"sqrt ( 4+5 )^2", 9}, {"sqrt(-0)", 0},
		{"50%", 0.5}, {"25%", 0.25}, {"2%", 0.02},
		{"200+10%", 220}, {"200-10%", 180}, {"200*10%", 20}, {"200/10%", 2000},
		{"200+10%+5%", 231}, {"25%+1", 1.25}, {"sqrt(25)*20%", 1},
		{"200+10%*2", 200.2}, {"200+-10%", 180}, {"200--10%", 220}, {"200+0.5%", 201},
		{"50%%", 0.005}, {"(50%)%", 0.005}, {"200+10%%", 200.2},
		{"200+(10%)", 220}, {"200+(-10%)", 180}, {"200+(5+5)%", 220},
		{"200+(10%+5%)", 200.105}, {"200+(10%*2)", 200.2},
		{"2^50%", math.Sqrt(2)}, {"200+2^50%", 200 + math.Sqrt(2)}, {"200+(2^3)%", 216},
		{"sqrt(100%)", 1}, {"200+sqrt(25)%", 210}, {"200+sqrt(100%)", 201},
		{strings.Repeat("(", MaxNestingDepth) + "1" + strings.Repeat(")", MaxNestingDepth), 1},
		{strings.Repeat("sqrt(", MaxNestingDepth) + "1" + strings.Repeat(")", MaxNestingDepth), 1},
		{strings.Repeat("1^", MaxNestingDepth) + "1", 1},
		{strings.Repeat("1+", 511) + "1", 512},
	} {
		t.Run(tc.expression, func(t *testing.T) {
			got, err := Evaluate(tc.expression)
			if err != nil || !closeEnough(got, tc.want) {
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
		{"1..2", "invalid_expression"}, {"1 2", "invalid_expression"}, {"2&3", "invalid_expression"},
		{"NaN", "invalid_expression"}, {"Infinity", "invalid_expression"}, {"cos(0)", "invalid_expression"},
		{"1e", "invalid_expression"}, {"1e+", "invalid_expression"}, {"0x10", "invalid_expression"},
		{"%2", "invalid_expression"}, {"1e999", "numeric_overflow"}, {strings.Repeat("9", 400), "numeric_overflow"},
		{"1e308*10", "numeric_overflow"}, {"1e308+1e308", "numeric_overflow"},
		{"-1e308-1e308", "numeric_overflow"}, {"1e308/1e-308", "numeric_overflow"},
		{strings.Repeat("1", 1025), "expression_too_large"},
		{"2^", "invalid_expression"}, {"2^^3", "invalid_expression"}, {"2^*3", "invalid_expression"},
		{"0^-1", "division_by_zero"}, {"(-2)^0.5", "domain_error"}, {"10^400", "numeric_overflow"},
		{"sqrt(-1)", "domain_error"}, {"sqrt()", "invalid_expression"}, {"sqrt", "invalid_expression"},
		{"sqrt(9", "invalid_expression"}, {"sqrt 9", "invalid_expression"}, {"sqrtx(9)", "invalid_expression"},
		{"()", "invalid_expression"}, {"(2+3", "invalid_expression"}, {"2+3)", "invalid_expression"},
		{"2(3)", "invalid_expression"}, {"sqrt(4)2", "invalid_expression"}, {"2%3", "invalid_expression"},
		{"200/0%", "division_by_zero"}, {"1/(1-1)%", "division_by_zero"},
		{"1e308+200%", "numeric_overflow"},
		{strings.Repeat("(", MaxNestingDepth+1) + "1" + strings.Repeat(")", MaxNestingDepth+1), "expression_too_complex"},
		{strings.Repeat("1^", MaxNestingDepth+1) + "1", "expression_too_complex"},
		{"(" + strings.Repeat("1^", MaxNestingDepth) + "1)", "expression_too_complex"},
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
	for _, s := range []string{"1+2*3", "5/0", "--.5", "1e308*2", "", "\xff", "2^3^2", "200+(10%)", "sqrt(-1)", "50%%", "(((1)))"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		value, err := Evaluate(s)
		if err == nil && (math.IsNaN(value) || math.IsInf(value, 0)) {
			t.Fatal("non-finite success")
		}
	})
}
