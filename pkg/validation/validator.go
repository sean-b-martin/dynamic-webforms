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
	datatypeRepository        *DatatypeRepository
	basicConstraintsValidator BasicConstraintsValidator
}

func NewFormValidator(repository *DatatypeRepository, basicConstraintsValidator BasicConstraintsValidator) *FormValidator {
	return &FormValidator{datatypeRepository: repository, basicConstraintsValidator: basicConstraintsValidator}
}

func (f *FormValidator) ValidateSchema(schema *model.WebFormSchema) []FormValidationError {
	helper := NewFormValidationHelper()

	// return if duplicate ids in schema
	if err := helper.ParseForm(schema, 5); err != nil {
		return err
	}

	var schemaErrors []FormValidationError

	for _, section := range schema.Sections {
		schemaErrors = append(schemaErrors, f.validateSchemaSection(section)...)
	}

	return schemaErrors
}

func (f *FormValidator) validateSchemaSection(schema *model.WebFormSection) []FormValidationError {
	var schemaErrors []FormValidationError
	for _, subsection := range schema.Subsections {
		schemaErrors = append(schemaErrors, f.validateSchemaSection(subsection)...)
	}

	for _, field := range schema.Fields {
		datatype, err := f.datatypeRepository.GetDatatype(field.Type)
		if err != nil {
			schemaErrors = append(schemaErrors, NewSchemaError(field.ID, "datatype does not exist"))
			continue
		}

		schemaErrors = append(schemaErrors, f.basicConstraintsValidator.ValidateSchema(field)...)
		schemaErrors = append(schemaErrors, datatype.ValidateSchema(field)...)
	}

	return schemaErrors
}

func (f *FormValidator) ValidateData(values []*model.WebFormDataRaw, schema *model.WebFormSchema) []FormValidationError {
	helper := NewFormValidationHelper()

	// stop validation if duplicate id's exist
	// only possible if schema was not validated using FormValidator.ValidateSchema(schema)
	if err := helper.ParseForm(schema, 5); err != nil {
		errs := make([]FormValidationError, 0, len(err)+1)
		errs = append(errs, NewSchemaError(-1, "form schema is invalid"))
		errs = append(errs, err...)
		return errs
	}

	var validationErrors []FormValidationError

	for _, value := range values {
		subfield, err := helper.GetSubfield(value.SchemaElementID)

		if errors.Is(err, ErrWrongElementType) {
			validationErrors = append(validationErrors, NewDataError(value.SchemaElementID, "element of schema id is not a (sub)field"))
			continue
		} else if errors.Is(err, ErrElementNotFound) {
			validationErrors = append(validationErrors, NewDataError(value.SchemaElementID, "schema id does not exist"))
			continue
		}

		datatype, err := f.datatypeRepository.GetDatatype(subfield.Type)
		// only possible if schema was not validated using FormValidator.ValidateSchema(schema)
		if err != nil {
			validationErrors = append(validationErrors, NewSchemaError(value.SchemaElementID, "datatype does not exist"))
			continue
		}

		// basic validation
		validationErrors = append(validationErrors, f.basicConstraintsValidator.ValidateData(value, &subfield.ValidationSchema.BasicConstraints)...)

		// dynamic validation
		validationErrors = append(validationErrors, datatype.ValidateData(value, subfield.ValidationSchema)...)

	}

	return validationErrors
}

type AllowsSubfieldsValidator struct {
	AllowsSubfield bool
}

func (a *AllowsSubfieldsValidator) Validate(schema *model.WebFormField) []FormValidationError {
	if !a.AllowsSubfield && (schema.Subfields != nil || len(schema.Subfields) > 0) {
		return []FormValidationError{NewSchemaError(schema.ID, "subfields are not allowed to be set in schema")}
	}

	return nil
}
