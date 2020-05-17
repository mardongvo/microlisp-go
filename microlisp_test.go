package microlisp

import (
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
			Statement{
				SType: STExpression,
				Expression: []Statement{
					Statement{SType: STString, ValueString: "somef"},
					Statement{SType: STString, ValueString: "param1"},
					Statement{SType: STString, ValueString: "param2"},
					Statement{SType: STInt, ValueInt: 3},
					Statement{SType: STFloat, ValueFloat: 4.0},
					Statement{SType: STBool, ValueBool: true},
					Statement{SType: STBool, ValueBool: false},
				},
			},
		},
		{" ( func1 (func2 a b) c d (func3 e) f) ",
			nil,
			Statement{
				SType: STExpression,
				Expression: []Statement{
					Statement{SType: STString, ValueString: "func1"},
					Statement{
						SType: STExpression,
						Expression: []Statement{
							Statement{SType: STString, ValueString: "func2"},
							Statement{SType: STString, ValueString: "a"},
							Statement{SType: STString, ValueString: "b"},
						},
					},
					Statement{SType: STString, ValueString: "c"},
					Statement{SType: STString, ValueString: "d"},
					Statement{
						SType: STExpression,
						Expression: []Statement{
							Statement{SType: STString, ValueString: "func3"},
							Statement{SType: STString, ValueString: "e"},
						},
					},
					Statement{SType: STString, ValueString: "f"},
				},
			},
		},
		{"somef",
			nil,
			Statement{
				SType:       STString,
				ValueString: "somef",
			},
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
