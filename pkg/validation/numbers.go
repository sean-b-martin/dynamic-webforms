package validation

import (
	"encoding/json"
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"math/big"
)

var DefaultIntNumberType = Datatype{
	definition: DatatypeDefinition{
		Identifier:      "@default/int64",
		DisplayName:     "int64",
		AllowsSubfields: false,
		InheritsFrom:    "",
	},
	DatatypeValidator: numberTypeValidator[int64]{},
}

type Numbers interface {
	int64 | int | float64
}

type LargeNumbers interface {
	big.Int | big.Float
}

type DynamicNumberValidation[T Numbers | LargeNumbers] struct {
	Lt        *T   `json:"lt,omitempty"`
	Gt        *T   `json:"gt,omitempty"`
	Lte       *T   `json:"lte,omitempty"`
	Gte       *T   `json:"gte,omitempty"`
	MinDigits *int `json:"minDigits,omitempty"`
	MaxDigits *int `json:"maxDigits,omitempty"`
}

type numberTypeValidator[T Numbers] struct {
}

func (v numberTypeValidator[T]) ValidateSchema(schema *model.WebFormField) []FormValidationError {
	errors := make([]FormValidationError, 0)
	var dynamicConstraints DynamicNumberValidation[T]

	err := json.Unmarshal(schema.ValidationSchema.DynamicConstraints, &dynamicConstraints)
	if err != nil {
		errors = append(errors, NewSchemaError(schema.ID, err.Error()))
	}

	// TODO add validation for values, for example Lt must be < than Gt

	return errors
}

func (v numberTypeValidator[T]) ValidateData(data *model.WebFormDataRaw, id int, schema *model.WebFormValidationSchema) []FormValidationError {
	var dynamicConstraints DynamicNumberValidation[T]

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

	if values != nil {
		return schemaErrors
	}

	validationErrors := make([]FormValidationError, 0)

	for i, value := range values {
		if dynamicConstraints.Lt != nil && *dynamicConstraints.Lt >= value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not less than "+fmt.Sprintf("%v", *dynamicConstraints.Lt)))
		}

		if dynamicConstraints.Gt != nil && *dynamicConstraints.Gt <= value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not greater than "+fmt.Sprintf("%v", *dynamicConstraints.Gt)))
		}

		if dynamicConstraints.Lte != nil && *dynamicConstraints.Lte <= value {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value not less or equal than "+fmt.Sprintf("%v", *dynamicConstraints.Lte)))
		}

		if dynamicConstraints.Gte != nil && *dynamicConstraints.Gte <= value {
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
