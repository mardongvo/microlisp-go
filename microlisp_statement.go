package microlisp

import (
	"strconv"
)

// Create new statement from string token (common during parsing)
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

func NewStringStatement(inp string) Statement {
	return Statement{SType: STString, ValueString: inp}
}

func NewErrorStatement(inp string) Statement {
	return Statement{SType: STError, ValueError: inp}
}

func NewIntStatement(inp int) Statement {
	return Statement{SType: STInt, ValueInt: inp}
}

func NewFloatStatement(inp float32) Statement {
	return Statement{SType: STFloat, ValueFloat: inp}
}

func NewBoolStatement(inp bool) Statement {
	return Statement{SType: STBool, ValueBool: inp}
}

func NewFuzzyStatement(inp FuzzySetType) Statement {
	return Statement{SType: STFuzzy, ValueFuzzy: inp}
}

// Check statement equality
// TODO: add check FuzzySet values
func IsEqualStatements(s1 Statement, s2 Statement) bool {
	if s1.SType != s2.SType {
		return false
	}
	if s1.SType == STString {
		return s1.ValueString == s2.ValueString
	}
	if s1.SType == STInt {
		return s1.ValueInt == s2.ValueInt
	}
	if s1.SType == STFloat {
		return s1.ValueFloat == s2.ValueFloat
	}
	if s1.SType == STBool {
		return s1.ValueBool == s2.ValueBool
	}
	if s1.SType == STExpression {
		if len(s1.Expression) != len(s2.Expression) {
			return false
		}
		for i := range s1.Expression {
			if !IsEqualStatements(s1.Expression[i], s2.Expression[i]) {
				return false
			}
		}
		return true
	}
	if s1.SType == STError {
		return s1.ValueError == s2.ValueError
	}
	return false
}
