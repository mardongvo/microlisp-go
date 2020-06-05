package microlisp

import "fmt"

// Standart logic functions (first-order logic) with lazy evaluation
var StandartLogicFunctions = FunctionMap{
	"not": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) != 1 {
			return NewErrorStatement(fmt.Errorf("Function `not' required one param"))
		}
		v := Eval(funcs, env, &expr[0])
		if v.Type() == STError {
			return v
		}
		if v.Type() != STBool {
			return NewErrorStatement(fmt.Errorf("Function `not' expect bool param"))
		}
		return NewBoolStatement(!v.ValueBool())
	},
	"and": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) == 0 {
			return NewErrorStatement(fmt.Errorf("Function `and' required at least one param"))
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.Type() == STError {
				return v
			}
			if v.Type() != STBool {
				return NewErrorStatement(fmt.Errorf("Function `and' expect bool param"))
			}
			if !v.ValueBool() {
				return NewBoolStatement(false)
			}
		}
		return NewBoolStatement(true)
	},
	"or": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) == 0 {
			return NewErrorStatement(fmt.Errorf("Function `or' required at least one param"))
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.Type() == STError {
				return v
			}
			if v.Type() != STBool {
				return NewErrorStatement(fmt.Errorf("Function `or' expect bool param"))
			}
			if v.ValueBool() {
				return NewBoolStatement(true)
			}
		}
		return NewBoolStatement(false)
	},
	"if": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) != 3 {
			return NewErrorStatement(fmt.Errorf("Function `if' required 3 param"))
		}
		cond := Eval(funcs, env, &expr[0])
		if cond.Type() == STError {
			return cond
		}
		if cond.Type() != STBool {
			return NewErrorStatement(fmt.Errorf("Function `if' expect bool param in condition"))
		}
		if cond.ValueBool() {
			return Eval(funcs, env, &expr[1])
		} else {
			return Eval(funcs, env, &expr[2])
		}
	},
}
