package microlisp

import "fmt"

func NewFuzzySet(normalize bool, elems ...FuzzyElement) FuzzySetType {
	var res FuzzySetType = make(FuzzySetType, len(elems))
	var sum float32
	copy(res, elems)
	if normalize {
		for _, e := range res {
			sum += e.Percent
		}
		if sum != 0.0 {
			for i := range res {
				res[i].Percent = res[i].Percent / sum
			}
		}
	}
	return res
}

func FuzzyEq(set FuzzySetType, find Statement) Statement {
	for _, e := range set {
		if IsEqualStatements(e.Value, find) {
			return NewFloatStatement(e.Percent)
		}
	}
	return NewFloatStatement(0.0)
}

func FuzzyEqSlice(set FuzzySetType, find []Statement) Statement {
	var res float32
	for _, f := range find {
		res += FuzzyEq(set, f).ValueFloat()
	}
	return NewFloatStatement(res)
}

// Fuzzy logic functions (first-order logic)
var FuzzyLogicFunctions = FunctionMap{
	"fnot": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) != 1 {
			return NewErrorStatement(fmt.Errorf("Function `fnot' required one param"))
		}
		v := Eval(funcs, env, &expr[0])
		if v.Type() == STError {
			return v
		}
		if v.Type() != STFloat {
			return NewErrorStatement(fmt.Errorf("Function `fnot' expect float param"))
		}
		return NewFloatStatement(1.0 - v.ValueFloat())
	},
	"fand": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		var res float32 = 1.0
		if len(expr) == 0 {
			return NewErrorStatement(fmt.Errorf("Function `fand' required at least one param"))
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.Type() == STError {
				return v
			}
			if v.Type() != STFloat {
				return NewErrorStatement(fmt.Errorf("Function `fand' expect float param"))
			}
			if v.ValueFloat() < res {
				res = v.ValueFloat()
			}
		}
		return NewFloatStatement(res)
	},
	"for": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		var res float32 = 0.0
		if len(expr) == 0 {
			return NewErrorStatement(fmt.Errorf("Function `for' required at least one param"))
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.Type() == STError {
				return v
			}
			if v.Type() != STFloat {
				return NewErrorStatement(fmt.Errorf("Function `for' expect float param"))
			}
			if v.ValueFloat() > res {
				res = v.ValueFloat()
			}
		}
		return NewFloatStatement(res)
	},
	/*
		// TODO: make correct fuzzy ternary op
		"fif": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
			if len(expr) != 3 {
				return NewErrorStatement("Function `if' required 3 param")
			}
			cond := Eval(funcs, env, &expr[0])
			if cond.SType == STError {
				return cond
			}
			if cond.SType != STBool {
				return NewErrorStatement("Function `if' expect bool param in condition")
			}
			if cond.ValueBool {
				return Eval(funcs, env, &expr[1])
			} else {
				return Eval(funcs, env, &expr[2])
			}
		},
	*/
}
