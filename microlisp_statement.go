package microlisp

import (
	"strconv"
)

// Create new statement from string token (common during parsing)
func NewStatement(inp string, tryConvert bool) Statement {
	if !tryConvert {
		return Statement{SType: STString, Value: inp}
	}
	if inp == "true" {
		return Statement{SType: STBool, Value: true}
	}
	if inp == "false" {
		return Statement{SType: STBool, Value: false}
	}
	i, err := strconv.ParseInt(inp, 10, 32)
	if err == nil {
		return Statement{SType: STInt, Value: int(i)}
	}
	f, err := strconv.ParseFloat(inp, 32)
	if err == nil {
		return Statement{SType: STFloat, Value: float32(f)}
	}
	return Statement{SType: STString, Value: inp}
}

func NewExpressionStatement(inp []Statement) Statement {
	return Statement{SType: STExpression, Value: inp}
}

func NewStringStatement(inp string) Statement {
	return Statement{SType: STString, Value: inp}
}

func NewErrorStatement(inp string) Statement {
	return Statement{SType: STError, Value: inp}
}

func NewIntStatement(inp int) Statement {
	return Statement{SType: STInt, Value: inp}
}

func NewFloatStatement(inp float32) Statement {
	return Statement{SType: STFloat, Value: inp}
}

func NewBoolStatement(inp bool) Statement {
	return Statement{SType: STBool, Value: inp}
}

func NewFuzzyStatement(inp FuzzySetType) Statement {
	return Statement{SType: STFuzzy, Value: inp}
}

//return
func (s Statement) ValueExpression() []Statement {
	if s.SType == STExpression {
		return s.Value.([]Statement)
	}
	return make([]Statement, 0)
}

func (s Statement) ValueString() string {
	if s.SType == STString {
		return s.Value.(string)
	}
	return ""
}

func (s Statement) ValueError() string {
	if s.SType == STError {
		return s.Value.(string)
	}
	return ""
}

func (s Statement) ValueInt() int {
	if s.SType == STInt {
		return s.Value.(int)
	}
	return 0
}

func (s Statement) ValueFloat() float32 {
	if s.SType == STFloat {
		return s.Value.(float32)
	}
	if s.SType == STInt {
		return float32(s.Value.(int))
	}
	return 0.0
}

func (s Statement) ValueBool() bool {
	if s.SType == STBool {
		return s.Value.(bool)
	}
	return false
}

// Check statement equality
// TODO: add check FuzzySet values
func IsEqualStatements(s1 Statement, s2 Statement) bool {
	if s1.SType != s2.SType {
		return false
	}
	if s1.SType == STString {
		return s1.ValueString() == s2.ValueString()
	}
	if s1.SType == STInt {
		return s1.ValueInt() == s2.ValueInt()
	}
	if s1.SType == STFloat {
		return s1.ValueFloat() == s2.ValueFloat()
	}
	if s1.SType == STBool {
		return s1.ValueBool() == s2.ValueBool()
	}
	if s1.SType == STExpression {
		exp1 := s1.ValueExpression()
		exp2 := s2.ValueExpression()
		if len(exp1) != len(exp2) {
			return false
		}
		for i := range exp1 {
			if !IsEqualStatements(exp1[i], exp2[i]) {
				return false
			}
		}
		return true
	}
	if s1.SType == STError {
		return s1.ValueError() == s2.ValueError()
	}
	return false
}
