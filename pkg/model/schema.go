package model

import "encoding/json"

type Schema struct {
	Type       string             `json:"type"`
	Fields     map[string]*Schema `json:"fields,omitempty"`
	Validators []ValidatorSchema  `json:"validator,omitempty"`
}

type ValidatorSchema struct {
	Rule       string          `json:"rule"`
	Parameters json.RawMessage `json:"-"`
}
