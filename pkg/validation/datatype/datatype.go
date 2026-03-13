package datatype

import "encoding/json"

type Datatype interface {
	Parse(raw json.RawMessage, isList bool) (any, error)
}

type GenericDatatype[T any] struct {
}

func (GenericDatatype[T]) Parse(raw json.RawMessage) (any, error) {
	var result T
	err := json.Unmarshal(raw, &result)
	return result, err
}
