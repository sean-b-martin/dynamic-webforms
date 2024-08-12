package validation

import (
	"encoding/json"
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"strings"
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
)

type Number interface {
	int64 | int | float64
}

type DynamicNumberValidationSchema[T Number] struct {
	Lt        *T   `json:"lt,omitempty"`
	Gt        *T   `json:"gt,omitempty"`
	Lte       *T   `json:"lte,omitempty"`
	Gte       *T   `json:"gte,omitempty"`
	MinDigits *int `json:"minDigits,omitempty"`
	MaxDigits *int `json:"maxDigits,omitempty"`
}

type numberTypeValidator[T Number] struct {
	allowsSubfieldsValidator AllowsSubfieldsValidator
}

func (v numberTypeValidator[T]) ValidateSchema(schema *model.WebFormField) []FormValidationError {
	errors := make([]FormValidationError, 0)

	var dynamicConstraints DynamicNumberValidationSchema[T]
	if schema.ValidationSchema.DynamicConstraints != nil {
		err := json.Unmarshal(schema.ValidationSchema.DynamicConstraints, &dynamicConstraints)
		if err != nil {
			errors = append(errors, NewSchemaError(schema.ID, err.Error()))
		}
	}

	errors = append(errors, v.allowsSubfieldsValidator.Validate(schema)...)

	if dynamicConstraints.Lt != nil {
		if dynamicConstraints.Gt != nil && *dynamicConstraints.Lt <= *dynamicConstraints.Gt {
			errors = append(errors, NewDataError(schema.ID, "lt must be greater than gt"))
		}
		if dynamicConstraints.Gte != nil && *dynamicConstraints.Lt <= *dynamicConstraints.Gte {
			errors = append(errors, NewDataError(schema.ID, "lt must be greater than gte"))
		}
	}

	if dynamicConstraints.Lte != nil {
		if dynamicConstraints.Gt != nil && *dynamicConstraints.Lte <= *dynamicConstraints.Gt {
			errors = append(errors, NewDataError(schema.ID, "lte must be greater than gt"))
		}
		if dynamicConstraints.Gte != nil && *dynamicConstraints.Lte < *dynamicConstraints.Gte {
			errors = append(errors, NewDataError(schema.ID, "lte must be greater than or equal gte"))
		}
	}

	if dynamicConstraints.MinDigits != nil {
		if dynamicConstraints.MaxDigits != nil && *dynamicConstraints.MaxDigits < *dynamicConstraints.MinDigits {
			errors = append(errors, NewDataError(schema.ID, "minDigits must be greater than maxDigits"))
		}
	}

	if dynamicConstraints.MaxDigits != nil && *dynamicConstraints.MaxDigits <= 0 {
		errors = append(errors, NewDataError(schema.ID, "maxDigits must be greater than 0"))
	}

	return errors
}

func (v numberTypeValidator[T]) ValidateData(data *model.WebFormDataRaw, schema *model.WebFormValidationSchema) []FormValidationError {
	if schema.DynamicConstraints == nil || len(schema.DynamicConstraints) == 0 {
		return []FormValidationError{}
	}

	var dynamicConstraints DynamicNumberValidationSchema[T]
	id := data.SchemaElementID

	if err := json.Unmarshal(schema.DynamicConstraints, &dynamicConstraints); err != nil {
		return []FormValidationError{NewSchemaError(id, err.Error())}
	}

	values := make([]T, len(data.Data))
	schemaErrors := make([]FormValidationError, 0)

	for i, v := range data.Data {
		var value T
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

		valueStr := strings.Split(fmt.Sprintf("%v", value), ".")[0]
		if dynamicConstraints.MinDigits != nil && *dynamicConstraints.MinDigits > len(valueStr) {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value has less digits than "+fmt.Sprintf("%v", *dynamicConstraints.MinDigits)))
		}

		if dynamicConstraints.MaxDigits != nil && *dynamicConstraints.MaxDigits < len(valueStr) {
			validationErrors = append(validationErrors, NewDataErrorWithIndex(id, i, "value has more digits than "+fmt.Sprintf("%v", *dynamicConstraints.MinDigits)))
		}
	}

	return validationErrors
}
