package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"strconv"
)

type BasicConstraintsValidator struct {
	defaultMaxItems int
}

func NewBasicConstraintsValidator(defaultMaxItems int) BasicConstraintsValidator {
	return BasicConstraintsValidator{defaultMaxItems: defaultMaxItems}
}

func (b BasicConstraintsValidator) ValidateBasicConstraintsSchema(schema *model.WebFormField) []FormValidationError {
	errs := make([]FormValidationError, 0)

	fields := make([]*model.WebFormSubfield, 0, 1+len(schema.Subfields))
	fields = append(fields, schema.WebFormSubfield)
	fields = append(fields, schema.Subfields...)

	for _, field := range fields {
		if field.ValidationSchema.MaxItems != nil {
			if *field.ValidationSchema.MaxItems <= 0 {
				errs = append(errs, NewSchemaError(schema.ID, "max items must be greater than 0"))
			}

			if field.ValidationSchema.MinItems != nil && *field.ValidationSchema.MinItems > *field.ValidationSchema.MaxItems {
				errs = append(errs, NewSchemaError(schema.ID, "max items must be greater or equal to min items"))
			}
		}
	}

	return errs
}

func (b BasicConstraintsValidator) ValidateBasicConstraints(data *model.WebFormDataRaw, schema *model.BasicConstraints) []FormValidationError {
	return b.validateItemCount(data, schema)
}

func (b BasicConstraintsValidator) validateItemCount(data *model.WebFormDataRaw, schema *model.BasicConstraints) []FormValidationError {
	errs := make([]FormValidationError, 0)

	if schema.MinItems != nil && *schema.MinItems > len(data.Data) {
		errs = append(errs, NewDataError(data.SchemaElementID, "amount of data too low, minimum allowed "+strconv.Itoa(*schema.MinItems)))
	}

	maxItems := schema.MaxItems
	if maxItems == nil && schema.MinItems == nil {
		maxItems = &b.defaultMaxItems
	}

	if maxItems != nil && *maxItems < len(data.Data) {
		errs = append(errs, NewDataError(data.SchemaElementID, "amount of data too high, maximum allowed "+strconv.Itoa(*schema.MaxItems)))
	}

	return errs
}
