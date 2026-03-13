package model

import "encoding/json"

type Schema struct {
	ID         int                `json:"id"`
	Type       string             `json:"type"`
	Fields     map[string]*Schema `json:"fields,omitempty"`
	Validators []ValidatorSchema  `json:"validator,omitempty"`
}

type ValidatorSchema struct {
	Rule       string          `json:"rule"`
	Field      int             `json:"field"`
	Parameters json.RawMessage `json:"-"`
}
