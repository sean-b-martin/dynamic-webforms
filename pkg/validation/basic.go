package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
)

func ItemCount(data model.WebFormDataRaw, schema *model.BasicConstraints) []error {
	errs := make([]error, 0)

	if schema.MinItems != nil && *schema.MinItems > len(data.Data) {
		errs = append(errs, FormValidationError{
			Name:            "ItemCount",
			SchemaElementID: data.SchemaElementID,
			Message:         "amount of items too low",
		})
	}

	if schema.MinItems != nil && *schema.MaxItems < len(data.Data) {
		errs = append(errs, FormValidationError{
			Name:            "ItemCount",
			SchemaElementID: data.SchemaElementID,
			Message:         "amount of items too high",
		})
	}

	return errs
}
