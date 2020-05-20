package microlisp

// Standart logic functions (first-order logic) with lazy evaluation
var StandartLogicFunctions = FunctionMap{
	"not": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) != 1 {
			return NewErrorStatement("Function `not' required one param")
		}
		v := Eval(funcs, env, &expr[0])
		if v.SType == STError {
			return v
		}
		if v.SType != STBool {
			return NewErrorStatement("Function `not' expect bool param")
		}
		return NewBoolStatement(!v.ValueBool)
	},
	"and": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) == 0 {
			return NewErrorStatement("Function `and' required at least one param")
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.SType == STError {
				return v
			}
			if v.SType != STBool {
				return NewErrorStatement("Function `and' expect bool param")
			}
			if !v.ValueBool {
				return NewBoolStatement(false)
			}
		}
		return NewBoolStatement(true)
	},
	"or": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
		if len(expr) == 0 {
			return NewErrorStatement("Function `or' required at least one param")
		}
		for _, e := range expr {
			v := Eval(funcs, env, &e)
			if v.SType == STError {
				return v
			}
			if v.SType != STBool {
				return NewErrorStatement("Function `or' expect bool param")
			}
			if v.ValueBool {
				return NewBoolStatement(true)
			}
		}
		return NewBoolStatement(false)
	},
	"if": func(funcs *FunctionMap, env *Environment, expr []Statement) Statement {
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
}
