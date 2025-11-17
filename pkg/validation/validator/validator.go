package validator

import "encoding/json"

type Validator interface {
	Validate(message any) error
}

type Definition struct {
	Description string         `json:"description"`
	Attributes  map[string]any `json:"attributes"`
}

type Factory interface {
	NewValidator(attributes *json.RawMessage) (Validator, error)
	Info() Definition
}
