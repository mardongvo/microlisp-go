package microlisp

import (
	"fmt"
	"regexp"
	"strconv"
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
)

type Statement struct {
	SType       StatementType
	ValueString string
	ValueInt    int
	ValueFloat  float32
	ValueBool   bool
	ValueFuzzy  FuzzySetType
	Expression  []Statement
}

type FunctionMap map[string]FunctionHandler

type Environment map[string]Statement

type FunctionHandler func(funcs *FunctionMap, env *Environment, expr *Statement)

//

func NewEnvironment() Environment {
	return make(Environment)
}

//
func NewStatement(inp string, tryConvert bool) Statement {
	if !tryConvert {
		return Statement{SType: STString, ValueString: inp}
	}
	if inp == "true" {
		return Statement{SType: STBool, ValueBool: true}
	}
	if inp == "false" {
		return Statement{SType: STBool, ValueBool: false}
	}
	i, err := strconv.ParseInt(inp, 10, 32)
	if err == nil {
		return Statement{SType: STInt, ValueInt: int(i)}
	}
	f, err := strconv.ParseFloat(inp, 32)
	if err == nil {
		return Statement{SType: STFloat, ValueFloat: float32(f)}
	}
	return Statement{SType: STString, ValueString: inp}
}

//*Parse
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

func SplitToTokens(program string) (tokens Tokens) {
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
var EndOfExpressionError = fmt.Errorf("Unexprected end of expression")
var ExpectOpenError = fmt.Errorf("Expected opening parenthesis")

//we expect only
//1. one atom
//2. s-expression
func BuildAST(tokens Tokens, startpos int) (Statement, int, error) {
	var resStmnt Statement
	pos := startpos
	if (tokens[pos].typ == atomToken) && (len(tokens) == 1) {
		return NewStatement(tokens[pos].val, true), pos, nil
	}
	if tokens[pos].typ == openToken {
		resStmnt.SType = STExpression
		resStmnt.Expression = make([]Statement, 1)
		pos++
		if pos >= len(tokens) {
			return Statement{}, pos, EndOfExpressionError
		}
		isFirstToken := true
		for pos < len(tokens) && (tokens[pos].typ != closeToken) {
			if tokens[pos].typ == atomToken {
				resStmnt.Expression = append(resStmnt.Expression,
					NewStatement(tokens[pos].val, !isFirstToken)) // do not convert first token (function name)
			}
			if tokens[pos].typ == openToken { //function name may be s-expression that return string
				stm, newpos, err := BuildAST(tokens, pos)
				if err != nil {
					return Statement{}, newpos, err
				}
				resStmnt.Expression = append(resStmnt.Expression, stm)
				pos = newpos
			}
			isFirstToken = false
			pos++
		}
		if pos >= len(tokens) {
			return Statement{}, pos, EndOfExpressionError
		}
	} else {
		return Statement{}, pos, ExpectOpenError
	}
	return resStmnt, pos, nil
}
