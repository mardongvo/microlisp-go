package microlisp

//

func NewEnvironment() Environment {
	return make(Environment)
}

func (env Environment) Add(key string, val Statement) {
	env[key] = val
}

func (env Environment) Get(key string) (Statement, bool) {
	v, ok := env[key]
	return v, ok
}
