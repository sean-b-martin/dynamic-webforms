package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator"
)

type StringType struct {
	BaseValidatorDatatype[string]
}

func NewStringType() *StringType {
	return &StringType{BaseValidatorDatatype[string]{
		ValidatorRepository: NewValidatorRepository(WithValidators(map[string]validator.Factory{
			"min_length": nil,
		})),
		DatatypeDefinition: DatatypeDefinition{
			Name:        "string",
			Description: "datatype for strings",
		},
	},
	}
}
