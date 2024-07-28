package validation

import (
	"errors"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
)

type FormValidator struct {
}

func (f *FormValidator) ValidateSchema(schema *model.WebFormSchema) []error {
	return nil
}

func (f *FormValidator) ValidateData(values []*model.WebFormDataRaw, schema *model.WebFormSchema) []error {
	helper := NewFormValidationHelper()

	if err := helper.ParseForm(schema); err != nil {
		return []error{errors.New("dynamic-forms: form schema is invalid")}
	}

	return nil
}
