package microlisp

import (
	"fmt"
	"regexp"
)

//fuzzy logic support
type FuzzyElement struct {
	Value   Statement
	Percent float32
}

type FuzzySetType []FuzzyElement

//types declaration
type StatementType uint8

const (
	STExpression StatementType = iota
	STString
	STInt
	STFloat
	STBool
	STFuzzy
	STError
)

type Statement struct {
	SType StatementType
	Value interface{}
}

type FunctionMap map[string]FunctionHandler

type Environment map[string]Statement

type FunctionHandler func(funcs *FunctionMap, env *Environment, expr []Statement) Statement

//***Parse-->
//partially from https://github.com/veonik/go-lisp/blob/master/lisp/tokens.go
type Tokens []*Token

type tokenType uint8

type Token struct {
	typ tokenType
	val string
}

type Pattern struct {
	typ    tokenType
	regexp *regexp.Regexp
}

const (
	whitespaceToken tokenType = iota
	atomToken
	openToken
	closeToken
)

func patterns() []Pattern {
	return []Pattern{
		{whitespaceToken, regexp.MustCompile(`^\s+`)},
		{atomToken, regexp.MustCompile(`^([^\(\)\s]+)`)},
		{openToken, regexp.MustCompile(`^(\()`)},
		{closeToken, regexp.MustCompile(`^(\))`)},
	}
}

func splitToTokens(program string) (tokens Tokens) {
	for pos := 0; pos < len(program); {
		for _, pattern := range patterns() {
			if matches := pattern.regexp.FindStringSubmatch(program[pos:]); matches != nil {
				if (len(matches) > 1) && (pattern.typ != whitespaceToken) {
					tokens = append(tokens, &Token{pattern.typ, matches[1]})
				}
				pos = pos + len(matches[0])
				break
			}
		}
	}
	return
}

//*AST
var ErrorEndOfExpression = fmt.Errorf("Unexprected end of expression")
var ErrorExpectOpen = fmt.Errorf("Expected opening parenthesis")
var ErrorTooManyTokens = fmt.Errorf("Too many tokens")

//we expect only
//1. one atom
//2. s-expression
func buildAST(tokens Tokens, startpos int) (Statement, int, error) {
	var expression = make([]Statement, 0)
	pos := startpos
	if (tokens[pos].typ == atomToken) && (len(tokens) == 1) {
		return NewStatement(tokens[pos].val, true), pos, nil
	}
	if tokens[pos].typ == openToken {
		pos++
		if pos >= len(tokens) {
			return Statement{}, pos, ErrorEndOfExpression
		}
		isFirstToken := true
		for pos < len(tokens) && (tokens[pos].typ != closeToken) {
			if tokens[pos].typ == atomToken {
				expression = append(expression,
					NewStatement(tokens[pos].val, !isFirstToken)) // do not convert first token (function name)
			}
			if tokens[pos].typ == openToken { //function name may be s-expression that return string
				stm, newpos, err := buildAST(tokens, pos)
				if err != nil {
					return Statement{}, newpos, err
				}
				expression = append(expression, stm)
				pos = newpos
			}
			isFirstToken = false
			pos++
		}
		if pos >= len(tokens) {
			return Statement{}, pos, ErrorEndOfExpression
		}
	} else {
		return Statement{}, pos, ErrorExpectOpen
	}
	return NewExpressionStatement(expression), pos, nil
}

func Parse(program string) (Statement, error) {
	tokens := splitToTokens(program)
	stm, endpos, err := buildAST(tokens, 0)
	if err != nil {
		return Statement{}, err
	}
	if endpos != len(tokens)-1 {
		return Statement{}, ErrorTooManyTokens
	}
	return stm, nil
}

//***<--Parse

// `env' function
func GetFromEnv(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
	var key Statement
	if len(expr) != 1 {
		return NewErrorStatement("Function `env' expect 1 param")
	}
	if expr[0].SType == STExpression {
		key = Eval(funcs, env, &expr[0])
	} else {
		key = expr[0]
	}
	if key.SType != STString {
		return NewErrorStatement("Function `env' expect 1 param is string")
	}
	if val, ok := (*env).Get(key.ValueString()); ok {
		return val
	}
	return NewErrorStatement(fmt.Sprintf("Environment key `%s' not found", key.ValueString()))
}

//Eval
func Eval(funcs *FunctionMap, env *Environment, expr *Statement) Statement {
	if expr.SType == STExpression {
		e := expr.ValueExpression()
		if len(e) == 0 {
			return NewErrorStatement("Expression without function name")
		}
		if e[0].ValueString() == "env" {
			return GetFromEnv(funcs, env, e[1:])
		}
		if fhandler, ok := (*funcs)[e[0].ValueString()]; ok {
			return fhandler(funcs, env, e[1:])
		}
		return NewErrorStatement(fmt.Sprintf("Function %s not found", e[0].ValueString()))
	}
	return *expr
}
