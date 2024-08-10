package validation

import (
	"encoding/json"
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"math/big"
)

var (
	DefaultIntNumberType = Datatype{
		definition: DatatypeDefinition{
			Identifier:      "@default/int64",
			DisplayName:     "int64",
			AllowsSubfields: false,
			InheritsFrom:    "",
		},
		DatatypeValidator: numberTypeValidator[int64]{AllowsSubfieldsValidator{false}},
	}

	DefaultFloatNumberType = Datatype{
		definition: DatatypeDefinition{
			Identifier:      "@default/float64",
			DisplayName:     "float64",
			AllowsSubfields: false,
			InheritsFrom:    "",
		},
		DatatypeValidator: numberTypeValidator[float64]{AllowsSubfieldsValidator{false}},
	}

	DefaultBigIntNumberType = Datatype{
		definition: DatatypeDefinition{
			Identifier:      "@default/large_int",
			DisplayName:     "large int",
			AllowsSubfields: false,
			InheritsFrom:    "",
		},
		DatatypeValidator: bigNumberTypeValidator[big.Int]{AllowsSubfieldsValidator{false}},
	}

	DefaultBigFloatNumberType = Datatype{
		definition: DatatypeDefinition{
			Identifier:      "@default/large_float",
			DisplayName:     "large float",
			AllowsSubfields: false,
			InheritsFrom:    "",
		},
		DatatypeValidator: bigNumberTypeValidator[big.Float]{AllowsSubfieldsValidator{false}},
	}
)

type Numbers interface {
	int64 | int | float64
}

type BigNumbers interface {
	big.Int | big.Float
}

type DynamicNumberValidationSchema[T Numbers | BigNumbers] struct {
	Lt        *T   `json:"lt,omitempty"`
	Gt        *T   `json:"gt,omitempty"`
	Lte       *T   `json:"lte,omitempty"`
	Gte       *T   `json:"gte,omitempty"`
	MinDigits *int `json:"minDigits,omitempty"`
	MaxDigits *int `json:"maxDigits,omitempty"`
}

type numberTypeValidator[T Numbers] struct {
	AllowsSubfieldsValidator
}

func (v numberTypeValidator[T]) ValidateSchema(schema *model.WebFormField) []FormValidationError {
	errors := make([]FormValidationError, 0)

	if schema.ValidationSchema.DynamicConstraints != nil {
		var dynamicConstraints DynamicNumberValidationSchema[T]
		err := json.Unmarshal(schema.ValidationSchema.DynamicConstraints, &dynamicConstraints)
		if err != nil {
			errors = append(errors, NewSchemaError(schema.ID, err.Error()))
		}
	}

	errors = append(errors, v.AllowsSubfieldsValidator.Validate(schema)...)

	// TODO add validation for values, for example Lt must be < than Gt

	return errors
}

func (v numberTypeValidator[T]) ValidateData(data *model.WebFormDataRaw, schema *model.WebFormValidationSchema) []FormValidationError {
	var dynamicConstraints DynamicNumberValidationSchema[T]
	id := data.SchemaElementID

	if err := json.Unmarshal(schema.DynamicConstraints, &dynamicConstraints); err != nil {
		return []FormValidationError{NewSchemaError(id, err.Error())}
	}

	values := make([]T, len(data.Data))
	schemaErrors := make([]FormValidationError, 0)

	for i, v := range data.Data {
		var value T
		// TODO cast strings to T if possible so "1" is also valid
		if err := json.Unmarshal(v, &value); err != nil {
			schemaErrors = append(schemaErrors, NewSchemaErrorWithIndex(id, i, err.Error()))
		}
		values[i] = value
	}

	if len(schemaErrors) > 0 {
		return schemaErrors
	}

	validationErrors := make([]FormValidationError, 0)

	for i, value := range values {
		if dynamicConstraints.Lt != nil && *dynamicConstraints.Lt <= value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not less than "+fmt.Sprintf("%v", *dynamicConstraints.Lt)))
		}

		if dynamicConstraints.Gt != nil && *dynamicConstraints.Gt >= value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not greater than "+fmt.Sprintf("%v", *dynamicConstraints.Gt)))
		}

		if dynamicConstraints.Lte != nil && *dynamicConstraints.Lte < value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not less or equal than "+fmt.Sprintf("%v", *dynamicConstraints.Lte)))
		}

		if dynamicConstraints.Gte != nil && *dynamicConstraints.Gte > value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not greater or equal than "+fmt.Sprintf("%v", *dynamicConstraints.Gte)))
		}

		valueStr := fmt.Sprintf("%v", value)
		if dynamicConstraints.MinDigits != nil && *dynamicConstraints.MinDigits > len(valueStr) {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value has less digits than "+fmt.Sprintf("%v", *dynamicConstraints.MinDigits)))
		}

		if dynamicConstraints.MaxDigits != nil && *dynamicConstraints.MaxDigits < len(valueStr) {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value has more digits than "+fmt.Sprintf("%v", *dynamicConstraints.MinDigits)))
		}
	}

	return validationErrors
}

type bigNumberTypeValidator[T BigNumbers] struct {
	AllowsSubfieldsValidator
}

func (b bigNumberTypeValidator[T]) ValidateSchema(schema *model.WebFormField) []FormValidationError {
	errors := make([]FormValidationError, 0)

	if schema.ValidationSchema.DynamicConstraints != nil {
		var dynamicConstraints DynamicNumberValidationSchema[T]
		err := json.Unmarshal(schema.ValidationSchema.DynamicConstraints, &dynamicConstraints)
		if err != nil {
			errors = append(errors, NewSchemaError(schema.ID, err.Error()))
		}
	}

	errors = append(errors, b.AllowsSubfieldsValidator.Validate(schema)...)
	// TODO add validation for values, for example Lt must be < than Gt

	return errors
}

func (b bigNumberTypeValidator[T]) ValidateData(data *model.WebFormDataRaw, schema *model.WebFormValidationSchema) []FormValidationError {
	//TODO implement me
	panic("implement me")
}
