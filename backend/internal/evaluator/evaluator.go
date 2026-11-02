// Package evaluator evaluates a deliberately restricted arithmetic grammar.
package evaluator

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

const MaxExpressionBytes = 1024
const MaxNestingDepth = 64

// Error is safe to expose through the API.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

func invalid() *Error {
	return &Error{"invalid_expression", "Enter a valid expression using numbers, operators, parentheses, and sqrt(...)"}
}
func overflow() *Error {
	return &Error{"numeric_overflow", "A number or result is outside the supported range"}
}

// Evaluate uses float64 with contextual percentages and right-associative powers.
// expression = term { ('+'|'-') term }; term = unary { ('*'|'/') unary }.
// unary = { '+'|'-' } power; power = postfix [ '^' unary ]; postfix = primary { '%' }.
// primary = number | '(' expression ')' | 'sqrt(' expression ')'.
func Evaluate(expression string) (float64, error) {
	if len(expression) > MaxExpressionBytes {
		return 0, &Error{"expression_too_large", "Expression must be at most 1024 bytes"}
	}
	p := parser{input: expression}
	value, err := p.expression()
	p.space()
	if err != nil {
		return 0, err
	}
	if p.pos != len(p.input) {
		return 0, invalid()
	}
	return value.number, nil
}

type parser struct {
	input string
	pos   int
	depth int
}

// directPercent survives only grouping/unary signs, not binary operations or sqrt.
// It distinguishes A+B% from A+B%*C without a second evaluation pass or an AST.
type operand struct {
	number        float64
	directPercent bool
}

func finite(number float64) (operand, error) {
	if math.IsInf(number, 0) || math.IsNaN(number) {
		return operand{}, overflow()
	}
	return operand{number: number}, nil
}

func (p *parser) enter() error {
	if p.depth >= MaxNestingDepth {
		return &Error{"expression_too_complex", "Expression exceeds the nesting limit of 64"}
	}
	p.depth++
	return nil
}

func (p *parser) take(token byte) bool {
	p.space()
	if p.pos < len(p.input) && p.input[p.pos] == token {
		p.pos++
		return true
	}
	return false
}

func (p *parser) space() {
	rest := strings.TrimLeftFunc(p.input[p.pos:], unicode.IsSpace)
	p.pos = len(p.input) - len(rest)
}

func (p *parser) expression() (operand, error) {
	value, err := p.term()
	if err != nil {
		return operand{}, err
	}
	for {
		p.space()
		if p.pos == len(p.input) || (p.input[p.pos] != '+' && p.input[p.pos] != '-') {
			return value, nil
		}
		op := p.input[p.pos]
		p.pos++
		right, err := p.term()
		if err != nil {
			return operand{}, err
		}
		amount := right.number
		if right.directPercent {
			amount *= value.number
		}
		var result float64
		if op == '+' {
			result = value.number + amount
		} else {
			result = value.number - amount
		}
		value, err = finite(result)
		if err != nil {
			return operand{}, err
		}
	}
}

func (p *parser) term() (operand, error) {
	value, err := p.unary()
	if err != nil {
		return operand{}, err
	}
	for {
		p.space()
		if p.pos == len(p.input) || (p.input[p.pos] != '*' && p.input[p.pos] != '/') {
			return value, nil
		}
		op := p.input[p.pos]
		p.pos++
		right, err := p.unary()
		if err != nil {
			return operand{}, err
		}
		if op == '/' && right.number == 0 {
			return operand{}, &Error{"division_by_zero", "Cannot divide by zero"}
		}
		var result float64
		if op == '*' {
			result = value.number * right.number
		} else {
			result = value.number / right.number
		}
		value, err = finite(result)
		if err != nil {
			return operand{}, err
		}
	}
}

func (p *parser) unary() (operand, error) {
	p.space()
	negative := false
	for p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
		if p.input[p.pos] == '-' {
			negative = !negative
		}
		p.pos++
		p.space()
	}
	value, err := p.power()
	if negative {
		value.number = -value.number
	}
	return value, err
}

func (p *parser) power() (operand, error) {
	left, err := p.postfix()
	if err != nil || !p.take('^') {
		return left, err
	}
	if err := p.enter(); err != nil {
		return operand{}, err
	}
	defer func() { p.depth-- }()
	right, err := p.unary()
	if err != nil {
		return operand{}, err
	}
	if left.number == 0 && right.number < 0 {
		return operand{}, &Error{"division_by_zero", "Cannot raise zero to a negative power"}
	}
	result := math.Pow(left.number, right.number)
	if math.IsNaN(result) {
		return operand{}, &Error{"domain_error", "Power is outside the real-number domain"}
	}
	return finite(result)
}

func (p *parser) postfix() (operand, error) {
	value, err := p.primary()
	if err != nil {
		return operand{}, err
	}
	for p.take('%') {
		value.number /= 100
		value.directPercent = true
	}
	return value, nil
}

func (p *parser) primary() (operand, error) {
	p.space()
	root := strings.HasPrefix(p.input[p.pos:], "sqrt")
	if root {
		p.pos += len("sqrt")
	}
	if root || (p.pos < len(p.input) && p.input[p.pos] == '(') {
		if !p.take('(') {
			return operand{}, invalid()
		}
		if err := p.enter(); err != nil {
			return operand{}, err
		}
		defer func() { p.depth-- }()
		value, err := p.expression()
		if err != nil {
			return operand{}, err
		}
		if !p.take(')') {
			return operand{}, invalid()
		}
		if root {
			if value.number < 0 {
				return operand{}, &Error{"domain_error", "Cannot take the square root of a negative number"}
			}
			return finite(math.Sqrt(value.number))
		}
		return value, nil
	}
	value, err := p.number()
	return operand{number: value}, err
}

func digit(b byte) bool { return b >= '0' && b <= '9' }

func (p *parser) number() (float64, error) {
	start := p.pos
	digits := 0
	for p.pos < len(p.input) && digit(p.input[p.pos]) {
		p.pos++
		digits++
	}
	if p.pos < len(p.input) && p.input[p.pos] == '.' {
		p.pos++
		for p.pos < len(p.input) && digit(p.input[p.pos]) {
			p.pos++
			digits++
		}
	}
	if digits == 0 {
		return 0, invalid()
	}
	if p.pos < len(p.input) && (p.input[p.pos] == 'e' || p.input[p.pos] == 'E') {
		p.pos++
		if p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
			p.pos++
		}
		exponentStart := p.pos
		for p.pos < len(p.input) && digit(p.input[p.pos]) {
			p.pos++
		}
		if exponentStart == p.pos {
			return 0, invalid()
		}
	}
	value, err := strconv.ParseFloat(p.input[start:p.pos], 64)
	if err != nil || math.IsInf(value, 0) || math.IsNaN(value) {
		return 0, overflow()
	}
	return value, nil
}
