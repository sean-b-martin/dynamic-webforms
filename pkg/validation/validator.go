package validation

import (
	"errors"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
)

type DatatypeValidator interface {
	ValidateSchema(schema *model.WebFormField) []FormValidationError
	ValidateData(data *model.WebFormDataRaw, schema *model.WebFormValidationSchema) []FormValidationError
}

type FormValidator struct {
	datatypeRepository *DatatypeRepository
}

func NewFormValidator(repository *DatatypeRepository) *FormValidator {
	return &FormValidator{datatypeRepository: repository}
}

func (f *FormValidator) ValidateSchema(schema *model.WebFormSchema) []FormValidationError {
	helper := NewFormValidationHelper()

	// return if duplicate ids in schema
	if err := helper.ParseForm(schema); err != nil {
		return err
	}

	schemaErrors := make([]FormValidationError, 0)

	for _, section := range schema.Sections {
		schemaErrors = append(schemaErrors, f.validateSchemaSection(section)...)
	}

	return schemaErrors
}

func (f *FormValidator) validateSchemaSection(schema *model.WebFormSection) []FormValidationError {
	schemaErrors := make([]FormValidationError, 0)
	for _, subsection := range schema.Subsections {
		schemaErrors = append(schemaErrors, f.validateSchemaSection(subsection)...)
	}

	for _, field := range schema.Fields {
		datatype, err := f.datatypeRepository.GetDatatype(field.Type)
		if err != nil {
			schemaErrors = append(schemaErrors, NewSchemaError(field.ID, "datatype does not exist"))
			continue
		}

		schemaErrors = append(schemaErrors, datatype.ValidateSchema(field)...)
	}

	return schemaErrors
}

func (f *FormValidator) ValidateData(values []*model.WebFormDataRaw, schema *model.WebFormSchema) []FormValidationError {
	helper := NewFormValidationHelper()

	// stop validation if duplicate id's exist
	// only possible if schema was not validated using FormValidator.ValidateSchema(schema)
	if err := helper.ParseForm(schema); err != nil {
		errs := make([]FormValidationError, 0, len(err)+1)
		errs = append(errs, NewSchemaError(-1, "form schema is invalid"))
		errs = append(errs, err...)
		return errs
	}

	validationErrors := make([]FormValidationError, 0)

	for _, value := range values {
		subfield, err := helper.GetSubfield(value.SchemaElementID)

		if errors.Is(err, WrongElementTypeError) {
			validationErrors = append(validationErrors, NewDataError(value.SchemaElementID, "element of schema id is not a (sub)field"))
			continue
		} else if errors.Is(err, ElementNotFoundError) {
			validationErrors = append(validationErrors, NewDataError(value.SchemaElementID, "schema id does not exist"))
			continue
		}

		datatype, err := f.datatypeRepository.GetDatatype(subfield.Type)
		// only possible if schema was not validated using FormValidator.ValidateSchema(schema)
		if err != nil {
			validationErrors = append(validationErrors, NewSchemaError(value.SchemaElementID, "datatype does not exist"))
			continue
		}

		validationErrors = append(validationErrors, datatype.ValidateData(value, subfield.ValidationSchema)...)

	}

	return validationErrors
}
