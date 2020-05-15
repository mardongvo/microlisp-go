package microlisp

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
