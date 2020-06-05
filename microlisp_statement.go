package microlisp

import (
	"fmt"
	"strconv"
)

// Create new statement from string token (common during parsing)
func NewStatement(inp string, tryConvert bool) Statement {
	if !tryConvert {
		return Statement{Value: inp}
	}
	if inp == "true" {
		return Statement{Value: true}
	}
	if inp == "false" {
		return Statement{Value: false}
	}
	i, err := strconv.ParseInt(inp, 10, 32)
	if err == nil {
		return Statement{Value: int(i)}
	}
	f, err := strconv.ParseFloat(inp, 32)
	if err == nil {
		return Statement{Value: float32(f)}
	}
	return Statement{Value: inp}
}

func NewExpressionStatement(inp []Statement) Statement {
	return Statement{inp}
}

func NewStringStatement(inp string) Statement {
	return Statement{inp}
}

func NewErrorStatement(inp error) Statement {
	return Statement{inp}
}

func NewIntStatement(inp int) Statement {
	return Statement{inp}
}

func NewFloatStatement(inp float32) Statement {
	return Statement{inp}
}

func NewBoolStatement(inp bool) Statement {
	return Statement{inp}
}

func NewFuzzyStatement(inp FuzzySetType) Statement {
	return Statement{inp}
}

//return
func (s Statement) Type() StatementType {
	switch s.Value.(type) {
	case []Statement:
		return STExpression
	case string:
		return STString
	case int:
		return STInt
	case float32:
		return STFloat
	case bool:
		return STBool
	case FuzzySetType:
		return STFuzzy
	case error:
		return STError
	default:
		return STUnknown
	}
}

func (s Statement) ValueExpression() []Statement {
	if s.Type() == STExpression {
		return s.Value.([]Statement)
	}
	return make([]Statement, 0)
}

func (s Statement) ValueString() string {
	if s.Type() == STString {
		return s.Value.(string)
	}
	return ""
}

func (s Statement) ValueError() error {
	if s.Type() == STError {
		return s.Value.(error)
	}
	return fmt.Errorf("")
}

func (s Statement) ValueInt() int {
	if s.Type() == STInt {
		return s.Value.(int)
	}
	return 0
}

func (s Statement) ValueFloat() float32 {
	if s.Type() == STFloat {
		return s.Value.(float32)
	}
	if s.Type() == STInt {
		return float32(s.Value.(int))
	}
	return 0.0
}

func (s Statement) ValueBool() bool {
	if s.Type() == STBool {
		return s.Value.(bool)
	}
	return false
}

// Check statement equality
// TODO: add check FuzzySet values
func IsEqualStatements(s1 Statement, s2 Statement) bool {
	if s1.Type() != s2.Type() {
		return false
	}
	if s1.Type() == STString {
		return s1.ValueString() == s2.ValueString()
	}
	if s1.Type() == STInt {
		return s1.ValueInt() == s2.ValueInt()
	}
	if s1.Type() == STFloat {
		return s1.ValueFloat() == s2.ValueFloat()
	}
	if s1.Type() == STBool {
		return s1.ValueBool() == s2.ValueBool()
	}
	if s1.Type() == STExpression {
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
	if s1.Type() == STError {
		return s1.ValueError().Error() == s2.ValueError().Error()
	}
	return false
}
