package validator

import (
	"encoding/json"
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"strings"
)

type Number interface {
	IntegerNumber | FloatNumber
}

type IntegerNumber interface {
	~int32 | ~int64 | ~int
}

type FloatNumber interface {
	~float32 | ~float64
}

type DynamicNumberValidationSchema[T Number] struct {
	Lt        *T   `json:"lt,omitempty"`
	Gt        *T   `json:"gt,omitempty"`
	Lte       *T   `json:"lte,omitempty"`
	Gte       *T   `json:"gte,omitempty"`
	MinDigits *int `json:"minDigits,omitempty"`
	MaxDigits *int `json:"maxDigits,omitempty"`
}

type GenericNumberValidator[T Number] struct{}

func (validator *GenericNumberValidator[T]) ValidateData(data *model.WebFormDataRaw, rawConstraints json.RawMessage) FieldValidatorError {
	validatorErr := NewFieldValidatorError(data.SchemaElementID)

	var constraints DynamicNumberValidationSchema[T]
	if err := json.Unmarshal(rawConstraints, &constraints); err != nil {
		validatorErr.AddFailedConstraint(FailedConstraintError{
			Constraint: "schema",
			DataIndex:  -1,
			Message:    validatorErr.Error(),
			Config:     nil,
		})
		return validatorErr
	}

	for index, v := range data.Data {
		var value T
		if err := json.Unmarshal(v, &value); err != nil {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "datatype",
				DataIndex:  index,
				Message:    err.Error(),
				Config:     nil,
			})

			continue
		}

		if constraints.Lt != nil && *constraints.Lt <= value {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "lt",
				DataIndex:  index,
				Message:    "must be less than",
				Config:     *constraints.Lt,
			})
		}

		if constraints.Gt != nil && *constraints.Gt >= value {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "gt",
				DataIndex:  index,
				Message:    "must be greater than",
				Config:     *constraints.Gt,
			})
		}

		if constraints.Lte != nil && *constraints.Lte < value {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "lte",
				DataIndex:  index,
				Message:    "must be less than or equal",
				Config:     *constraints.Lte,
			})
		}

		if constraints.Gte != nil && *constraints.Gte > value {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "gte",
				DataIndex:  index,
				Message:    "must be greater than or equal",
				Config:     *constraints.Gte,
			})
		}

		if value < 0 {
			value = -value
		}

		valueStr := strings.Split(fmt.Sprintf("%v", value), ".")[0]
		if constraints.MinDigits != nil && *constraints.MinDigits > len(valueStr) {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "minDigits",
				DataIndex:  index,
				Message:    "must have more digits than",
				Config:     *constraints.MinDigits,
			})
		}

		if constraints.MaxDigits != nil && *constraints.MaxDigits < len(valueStr) {
			validatorErr.AddFailedConstraint(FailedConstraintError{
				Constraint: "maxDigits",
				DataIndex:  index,
				Message:    "must have less digits than",
				Config:     *constraints.MaxDigits,
			})
		}
	}

	return validatorErr
}
