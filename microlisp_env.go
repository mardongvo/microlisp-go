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

// Convert JSON object to Environment
// Subobjects interpret as FuzzySet
func JsonMapToEnvironment(inp map[string]interface{}) Environment {
	var res = NewEnvironment()
	for k, v := range inp {
		switch vv := v.(type) {
		case string:
			res.Add(k, NewStringStatement(vv))
		case int:
			res.Add(k, NewIntStatement(vv))
		case float64:
			res.Add(k, NewFloatStatement(float32(vv)))
		case bool:
			res.Add(k, NewBoolStatement(vv))
		case map[string]interface{}:
			{ //fuzzy set, expect {"stringkey": float...}
				fuz := make(FuzzySetType, 0)
				for k1, v1 := range vv {
					percent, ok := v1.(float64)
					if ok {
						fuz = append(fuz, FuzzyElement{NewStringStatement(k1), float32(percent)})
					}
				}
				res.Add(k, NewFuzzyStatement(fuz))
			}
		}
	}
	return res
}
