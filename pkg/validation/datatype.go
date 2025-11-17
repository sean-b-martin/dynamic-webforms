package validation

import (
	"encoding/json"
)

type Datatype interface {
	GetDefinition() DatatypeDefinition
	UnmarshalData(data *json.RawMessage, multiple bool) (any, error)
}

type DatatypeDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ValidatorDatatype interface {
	Datatype
	ValidatorRepository
}

type BaseValidatorDatatype[T any] struct {
	ValidatorRepository
	DatatypeDefinition
}

func (b *BaseValidatorDatatype[T]) UnmarshalData(data *json.RawMessage, asList bool) (any, error) {
	if asList {
		var value []T
		if err := json.Unmarshal(*data, &value); err != nil {
			return nil, err
		}
		return value, nil
	}

	var value T
	if err := json.Unmarshal(*data, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BaseValidatorDatatype[T]) GetDefinition() DatatypeDefinition {
	return DatatypeDefinition{b.DatatypeDefinition.Name, b.DatatypeDefinition.Description}
}
