package datatype

import "fmt"

type Registry struct {
	types map[string]Datatype
}

func NewRegistry() *Registry {
	return &Registry{
		types: make(map[string]Datatype),
	}
}

func (r *Registry) Register(name string, t Datatype) error {
	if _, ok := r.types[name]; ok {
		return fmt.Errorf("%s is already registered", name)
	}
	r.types[name] = t
	return nil
}

func (r *Registry) Get(name string) (Datatype, error) {
	t, ok := r.types[name]
	if !ok {
		return nil, fmt.Errorf("%s is not registered", name)
	}
	return t, nil
}
