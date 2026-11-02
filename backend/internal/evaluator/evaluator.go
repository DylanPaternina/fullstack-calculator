// Package evaluator evaluates a deliberately restricted arithmetic grammar.
package evaluator

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

const MaxExpressionBytes = 1024

// Error is safe to expose through the API.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

func invalid() *Error {
	return &Error{"invalid_expression", "Enter a valid expression using numbers and +, -, *, /"}
}
func overflow() *Error {
	return &Error{"numeric_overflow", "A number or result is outside the supported range"}
}

// Evaluate applies standard precedence and left associativity using float64.
// Grammar: expression = term { ('+'|'-') term }; term = number { ('*'|'/') number }.
// A number can have unary signs and a decimal/scientific literal, but no grouping.
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
	return value, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) space() {
	rest := strings.TrimLeftFunc(p.input[p.pos:], unicode.IsSpace)
	p.pos = len(p.input) - len(rest)
}

func (p *parser) expression() (float64, error) {
	value, err := p.term()
	if err != nil {
		return 0, err
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
			return 0, err
		}
		if op == '+' {
			value += right
		} else {
			value -= right
		}
		if math.IsInf(value, 0) || math.IsNaN(value) {
			return 0, overflow()
		}
	}
}

func (p *parser) term() (float64, error) {
	value, err := p.number()
	if err != nil {
		return 0, err
	}
	for {
		p.space()
		if p.pos == len(p.input) || (p.input[p.pos] != '*' && p.input[p.pos] != '/') {
			return value, nil
		}
		op := p.input[p.pos]
		p.pos++
		right, err := p.number()
		if err != nil {
			return 0, err
		}
		if op == '/' && right == 0 {
			return 0, &Error{"division_by_zero", "Cannot divide by zero"}
		}
		if op == '*' {
			value *= right
		} else {
			value /= right
		}
		if math.IsInf(value, 0) || math.IsNaN(value) {
			return 0, overflow()
		}
	}
}

func digit(b byte) bool { return b >= '0' && b <= '9' }

func (p *parser) number() (float64, error) {
	p.space()
	negative := false
	for p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
		if p.input[p.pos] == '-' {
			negative = !negative
		}
		p.pos++
		p.space()
	}
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
	if negative {
		value = -value
	}
	return value, nil
}
