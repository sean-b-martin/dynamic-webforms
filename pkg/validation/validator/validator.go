package validator

import "encoding/json"

type Validator interface {
	Validate(*Context, *FieldRef) error
}

type Factory func(message json.RawMessage) (Validator, error)
