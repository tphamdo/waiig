package object

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

func NewEnclosedEnvironment(out *Environment) *Environment {
	return &Environment{store: make(map[string]Object), outer: out}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if e.outer != nil && !ok {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
