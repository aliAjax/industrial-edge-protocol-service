package rules

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenKind int

const (
	TokenIdentifier TokenKind = iota
	TokenNumber
	TokenOperator
	TokenLParen
	TokenRParen
)

type Token struct {
	Kind TokenKind
	Text string
}

func Lex(input string) ([]Token, error) {
	fields := strings.Fields(input)
	out := []Token{}
	for _, f := range fields {
		kind := TokenIdentifier
		if _, err := strconv.ParseFloat(f, 64); err == nil {
			kind = TokenNumber
		}
		if f == ">" || f == "<" || f == ">=" || f == "<=" || f == "==" || f == "+" || f == "-" || f == "*" || f == "/" {
			kind = TokenOperator
		}
		if f == "(" {
			kind = TokenLParen
		}
		if f == ")" {
			kind = TokenRParen
		}
		out = append(out, Token{Kind: kind, Text: f})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	return out, nil
}

type AST interface {
	Eval(map[string]float64) (float64, error)
	String() string
}
type Literal struct{ Value float64 }

func (l Literal) Eval(map[string]float64) (float64, error) { return l.Value, nil }
func (l Literal) String() string                           { return strconv.FormatFloat(l.Value, 'f', -1, 64) }

type Variable struct{ Name string }

func (v Variable) Eval(values map[string]float64) (float64, error) {
	value, ok := values[v.Name]
	if !ok {
		return 0, fmt.Errorf("unknown variable %s", v.Name)
	}
	return value, nil
}
func (v Variable) String() string { return v.Name }

type Binary struct {
	Left     AST
	Operator string
	Right    AST
}

func (b Binary) Eval(values map[string]float64) (float64, error) {
	l, err := b.Left.Eval(values)
	if err != nil {
		return 0, err
	}
	r, err := b.Right.Eval(values)
	if err != nil {
		return 0, err
	}
	switch b.Operator {
	case "+":
		return l + r, nil
	case "-":
		return l - r, nil
	case "*":
		return l * r, nil
	case "/":
		if r == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return l / r, nil
	default:
		return 0, fmt.Errorf("operator %s", b.Operator)
	}
}
func (b Binary) String() string { return "(" + b.Left.String() + b.Operator + b.Right.String() + ")" }
func ParseArithmetic(input string) (AST, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 1 {
		if tokens[0].Kind == TokenNumber {
			v, _ := strconv.ParseFloat(tokens[0].Text, 64)
			return Literal{v}, nil
		}
		return Variable{tokens[0].Text}, nil
	}
	if len(tokens) != 3 || tokens[1].Kind != TokenOperator {
		return nil, fmt.Errorf("expected binary expression")
	}
	left, err := ParseArithmetic(tokens[0].Text)
	if err != nil {
		return nil, err
	}
	right, err := ParseArithmetic(tokens[2].Text)
	if err != nil {
		return nil, err
	}
	return Binary{Left: left, Operator: tokens[1].Text, Right: right}, nil
}
