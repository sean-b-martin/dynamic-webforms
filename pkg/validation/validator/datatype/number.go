package datatype

import (
	"encoding/json"
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator/common"
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

type GenericNumberValidator[T Number] struct {
	constraints DynamicNumberValidationSchema[T]
}

func (validator *GenericNumberValidator[T]) New() Validator {
	return &GenericNumberValidator[T]{}
}

func (validator *GenericNumberValidator[T]) Initialize(id int, rawConstraints *json.RawMessage) common.FieldValidatorError {
	validatorErr := common.NewFieldValidatorError(id)

	if err := common.UnmarshalValidationSchema(rawConstraints, &validator.constraints); err != nil {
		validatorErr.AddFailedConstraint(*err)
		return validatorErr
	}

	if validator.constraints.Lt != nil {
		if validator.constraints.Gt != nil && *validator.constraints.Lt <= *validator.constraints.Gt {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lt",
				DataIndex:  -1,
				Message:    "lt must be greater than gt",
				Config:     nil,
			})
		}
		if validator.constraints.Gte != nil && *validator.constraints.Lt <= *validator.constraints.Gte {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lt",
				DataIndex:  -1,
				Message:    "lt must be greater than gte",
				Config:     nil,
			})
		}
	}

	if validator.constraints.Lte != nil {
		if validator.constraints.Gt != nil && *validator.constraints.Lte <= *validator.constraints.Gt {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lt",
				DataIndex:  -1,
				Message:    "lte must be greater than gt",
				Config:     nil,
			})
		}
		if validator.constraints.Gte != nil && *validator.constraints.Lte < *validator.constraints.Gte {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lt",
				DataIndex:  -1,
				Message:    "lte must be greater than or equal gte",
				Config:     nil,
			})
		}
	}

	if validator.constraints.MinDigits != nil {
		if validator.constraints.MaxDigits != nil && *validator.constraints.MaxDigits < *validator.constraints.MinDigits {
			if validator.constraints.Gte != nil && *validator.constraints.Lte < *validator.constraints.Gte {
				validatorErr.AddFailedConstraint(common.FailedConstraintError{
					Constraint: "lt",
					DataIndex:  -1,
					Message:    "minDigits must be greater than maxDigits",
					Config:     nil,
				})
			}
		}
	}

	if validator.constraints.MaxDigits != nil && *validator.constraints.MaxDigits <= 0 {
		validatorErr.AddFailedConstraint(common.FailedConstraintError{
			Constraint: "lt",
			DataIndex:  -1,
			Message:    "maxDigits must be greater than 0",
			Config:     nil,
		})
	}

	return validatorErr
}

func (validator *GenericNumberValidator[T]) Validate(data *model.WebFormDataRaw) common.FieldValidatorError {
	validatorErr := common.NewFieldValidatorError(data.SchemaElementID)

	constraints := validator.constraints

	for index, v := range data.Data {
		var value T
		if err := json.Unmarshal(v, &value); err != nil {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "datatype",
				DataIndex:  index,
				Message:    err.Error(),
				Config:     nil,
			})

			continue
		}

		if constraints.Lt != nil && *constraints.Lt <= value {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lt",
				DataIndex:  index,
				Message:    "must be less than",
				Config:     *constraints.Lt,
			})
		}

		if constraints.Gt != nil && *constraints.Gt >= value {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "gt",
				DataIndex:  index,
				Message:    "must be greater than",
				Config:     *constraints.Gt,
			})
		}

		if constraints.Lte != nil && *constraints.Lte < value {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "lte",
				DataIndex:  index,
				Message:    "must be less than or equal",
				Config:     *constraints.Lte,
			})
		}

		if constraints.Gte != nil && *constraints.Gte > value {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
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
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "minDigits",
				DataIndex:  index,
				Message:    "must have more digits than",
				Config:     *constraints.MinDigits,
			})
		}

		if constraints.MaxDigits != nil && *constraints.MaxDigits < len(valueStr) {
			validatorErr.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "maxDigits",
				DataIndex:  index,
				Message:    "must have less digits than",
				Config:     *constraints.MaxDigits,
			})
		}
	}

	return validatorErr
}
