package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator/datatype"
)

// DatatypeDefinition contains metadata about a datatype such as the id for using the datatype
// and whether the datatype allows subfields or not.
type DatatypeDefinition struct {
	ID              string `json:"ID"`
	DisplayName     string `json:"displayName"`
	AllowsSubfields bool   `json:"allowsSubfields"`
}

// Datatype contains its definition and all rules for validation
type Datatype struct {
	definition DatatypeDefinition
	datatype.Validator
}

type DatatypeDynamicConstraints struct {
	InheritedConstraints *DatatypeDynamicConstraints `json:"inheritedConstraints,omitempty"`
	Constraints          map[string]interface{}      `json:"rules,omitempty"`
}

func NewDatatype(definition DatatypeDefinition, datatypeValidator datatype.Validator) Datatype {
	return Datatype{definition: definition, Validator: datatypeValidator}
}
