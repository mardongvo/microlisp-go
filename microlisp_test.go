package microlisp

import (
	"math"
	"reflect"
	"testing"
)

func TestTokens1(t *testing.T) {
	var tests = []struct {
		inp  string
		outp Tokens
	}{
		{" ( somef ) ", Tokens{{openToken, "("}, {atomToken, "somef"}, {closeToken, ")"}}},
		{"somef", Tokens{{atomToken, "somef"}}},
		{")((somef))", Tokens{{closeToken, ")"}, {openToken, "("},
			{openToken, "("}, {atomToken, "somef"}, {closeToken, ")"},
			{closeToken, ")"}}},
	}
	for _, test := range tests {
		x := splitToTokens(test.inp)
		if !reflect.DeepEqual(x, test.outp) {
			t.Errorf("SplitToTokens \"%v\" gives \"%v\", expected \"%v\"",
				test.inp, x, test.outp)
		}
	}
}

func TestAST(t *testing.T) {
	var tests = []struct {
		inp  string
		err  error
		outp Statement
	}{
		{" ( somef param1 param2 3 4.0 true false) ",
			nil,
			NewExpressionStatement([]Statement{
				NewStringStatement("somef"),
				NewStringStatement("param1"),
				NewStringStatement("param2"),
				NewIntStatement(3),
				NewFloatStatement(4.0),
				NewBoolStatement(true),
				NewBoolStatement(false),
			}),
		},
		{" ( func1 (func2 a b) c d (func3 e) f) ",
			nil,
			NewExpressionStatement([]Statement{
				NewStringStatement("func1"),
				NewExpressionStatement(
					[]Statement{
						NewStringStatement("func2"),
						NewStringStatement("a"),
						NewStringStatement("b"),
					}),
				NewStringStatement("c"),
				NewStringStatement("d"),
				NewExpressionStatement(
					[]Statement{
						NewStringStatement("func3"),
						NewStringStatement("e"),
					}),
				NewStringStatement("f"),
			}),
		},
		{"somef",
			nil,
			NewStringStatement("somef"),
		},
		{"somef erratom",
			ErrorExpectOpen,
			Statement{},
		},
		{")somef erratom",
			ErrorExpectOpen,
			Statement{},
		},
		{"(somef erratom",
			ErrorEndOfExpression,
			Statement{},
		},
		{"(somef erratom()",
			ErrorEndOfExpression,
			Statement{},
		},
		{"(somef erratom))",
			ErrorTooManyTokens,
			Statement{},
		},
		{"(somef erratom) misfit",
			ErrorTooManyTokens,
			Statement{},
		},
	}
	for _, test := range tests {
		ast, err := Parse(test.inp)
		if err != test.err {
			t.Errorf("BuildAST \"%v\" error is \"%v\", expected \"%v\"",
				test.inp, err, test.err)
		}
		if !reflect.DeepEqual(ast, test.outp) {
			t.Errorf("BuildAST \"%v\" gives \"%#v\", expected \"%#v\"",
				test.inp, ast, test.outp)
		}
	}

}

func TestEval(t *testing.T) {
	var tests = []struct {
		program string
		funcs   FunctionMap
		env     Environment
		result  Statement
	}{
		{"(env somekey)",
			FunctionMap{},
			Environment{"somekey": NewStringStatement("somevalue")},
			NewStringStatement("somevalue"),
		},
		{"(env nokey)",
			FunctionMap{},
			Environment{"somekey": NewStringStatement("somevalue")},
			NewErrorStatement("Environment key `nokey' not found"),
		},
	}
	for _, test := range tests {
		ast, _ := Parse(test.program)
		val := Eval(&test.funcs, &test.env, &ast)
		if !IsEqualStatements(val, test.result) {
			t.Errorf("Eval \"%v\" gives \"%#v\", expected \"%#v\"",
				test.program, val, test.result)
		}
	}

}

func TestEvalStandartLogic(t *testing.T) {
	var tests = []struct {
		program string
		funcs   FunctionMap
		env     Environment
		result  Statement
	}{
		{"(and (env a) (env b))",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(true), "b": NewBoolStatement(true)},
			NewBoolStatement(true),
		},
		{"(and (env a) (env b))",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(true), "b": NewBoolStatement(false)},
			NewBoolStatement(false),
		},
		{"(or (env a) (env b) (env c))",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(false), "b": NewBoolStatement(false), "c": NewBoolStatement(false)},
			NewBoolStatement(false),
		},
		{"(or (or (env a) (env d)) (env b) (env c))",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(false), "b": NewBoolStatement(false), "c": NewBoolStatement(false),
				"d": NewBoolStatement(true)},
			NewBoolStatement(true),
		},
		{"(and (env a) (env b))",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(true), "b": NewErrorStatement("Wow!")},
			NewErrorStatement("Wow!"),
		},
		{"(if (env a) b c)",
			StandartLogicFunctions,
			Environment{"a": NewBoolStatement(true), "b": NewStringStatement("Here")},
			NewStringStatement("b"),
		},
	}
	for _, test := range tests {
		ast, _ := Parse(test.program)
		val := Eval(&test.funcs, &test.env, &ast)
		if !IsEqualStatements(val, test.result) {
			t.Errorf("Eval(standart logic) \"%v\" gives \"%#v\", expected \"%#v\"",
				test.program, val, test.result)
		}
	}
}

func TestEvalFuzzyLogic(t *testing.T) {
	var tests = []struct {
		program string
		funcs   FunctionMap
		env     Environment
		result  Statement
	}{
		{"(fand (env a) (env b))",
			FuzzyLogicFunctions,
			Environment{"a": NewFloatStatement(1.0), "b": NewFloatStatement(1.0)},
			NewFloatStatement(1.0),
		},
		{"(fand (env a) (env b))",
			FuzzyLogicFunctions,
			Environment{"a": NewFloatStatement(0.9), "b": NewFloatStatement(0.5)},
			NewFloatStatement(0.5),
		},
		{"(for (env a) (env b) (env c))",
			FuzzyLogicFunctions,
			Environment{"a": NewFloatStatement(0.9), "b": NewFloatStatement(0.8), "c": NewFloatStatement(0.7)},
			NewFloatStatement(0.9),
		},
		{"(for (for (env a) (env d)) (env b) (env c))",
			FuzzyLogicFunctions,
			Environment{"a": NewFloatStatement(0.1), "b": NewFloatStatement(0.5), "c": NewFloatStatement(0.3),
				"d": NewFloatStatement(0.7)},
			NewFloatStatement(0.7),
		},
		{"(fand (env a) (env b))",
			FuzzyLogicFunctions,
			Environment{"a": NewFloatStatement(0.1), "b": NewErrorStatement("Wow!")},
			NewErrorStatement("Wow!"),
		},
	}
	for _, test := range tests {
		ast, _ := Parse(test.program)
		val := Eval(&test.funcs, &test.env, &ast)
		if !IsEqualStatements(val, test.result) {
			t.Errorf("Eval(fuzzy logic) \"%v\" gives \"%#v\", expected \"%#v\"",
				test.program, val, test.result)
		}
	}
}

func TestFuzzyEq(t *testing.T) {
	var tests = []struct {
		set    FuzzySetType
		find   []Statement
		result Statement
	}{
		{
			NewFuzzySet(false, FuzzyElement{NewStringStatement("a"), 0.1},
				FuzzyElement{NewStringStatement("b"), 0.3},
				FuzzyElement{NewStringStatement("c"), 0.6}),
			[]Statement{NewStringStatement("c")},
			NewFloatStatement(0.6),
		},
		{
			NewFuzzySet(false, FuzzyElement{NewStringStatement("a"), 0.1},
				FuzzyElement{NewStringStatement("b"), 0.3},
				FuzzyElement{NewStringStatement("c"), 0.6}),
			[]Statement{NewStringStatement("c"), NewStringStatement("b")},
			NewFloatStatement(0.9),
		},
	}
	for _, test := range tests {
		val := FuzzyEqSlice(test.set, test.find)
		val = NewFloatStatement(float32(math.Round(float64(val.ValueFloat())*1000.0) / 1000.0))
		if !IsEqualStatements(val, test.result) {
			t.Errorf("FuzzyEq gives \"%#v\", expected \"%#v\"",
				val, test.result)
		}
	}
}

func TestJson2Env(t *testing.T) {
	inp := map[string]interface{}{
		"a": "abc",
		"b": true,
		"c": 10,
		"d": 5.0,
		"e": map[string]interface{}{
			"x": 0.1,
			"y": 0.9,
		},
		"f": nil, //skipped
		"g": []interface{}{1, 2, 3, "u"},
	}
	out := Environment{
		"a": NewStringStatement("abc"),
		"b": NewBoolStatement(true),
		"c": NewIntStatement(10),
		"d": NewFloatStatement(5.0),
		"e": NewFuzzyStatement(NewFuzzySet(false,
			FuzzyElement{NewStringStatement("x"), 0.1},
			FuzzyElement{NewStringStatement("y"), 0.9},
		)),
	}
	env := JsonMapToEnvironment(inp)
	for k, v1 := range env {
		v2, ok := out[k]
		if !ok {
			t.Errorf("JsonMapToEnvironment expect key \"%#v\" but nothing",
				k)
		}
		if !reflect.DeepEqual(v1, v2) {
			t.Errorf("JsonMapToEnvironment key %s: got \"%#v\", expected \"%#v\"",
				k, v1, v2)
		}
	}
}
