package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator/common"
)

type schemaElementType int

const (
	SECTION schemaElementType = iota
	FIELD
	SUBFIELD
)

type schemaElement struct {
	Type    schemaElementType
	Element interface{}
}

type FormValidationHelper struct {
	elements     map[int]schemaElement
	errors       []common.ValidatorError
	maxRecursion int
}

func NewFormValidationHelper(maxRecursion int) *FormValidationHelper {
	return &FormValidationHelper{
		elements:     make(map[int]schemaElement),
		errors:       nil,
		maxRecursion: maxRecursion,
	}
}

func (f *FormValidationHelper) ParseForm(schema *model.WebFormSchema) []common.ValidatorError {
	for _, section := range schema.Sections {
		f.parseSection(section, f.maxRecursion)
	}

	return f.errors
}

func (f *FormValidationHelper) parseSection(schema *model.WebFormSection, maxRecursion int) {
	if _, ok := f.elements[schema.ID]; ok {
		err := common.NewFieldValidatorError(schema.ID)
		err.AddFailedConstraint(common.FailedConstraintError{
			Constraint: "unique id",
			Message:    "id must be unique",
		})
		f.errors = append(f.errors, err)
	}

	f.elements[schema.ID] = schemaElement{SECTION, schema}

	for _, subsection := range schema.Subsections {
		if maxRecursion == 0 {
			err := common.NewFieldValidatorError(schema.ID)
			err.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "subsections",
				Message:    "recursion must be less than",
				Config:     f.maxRecursion,
			})
			f.errors = append(f.errors, err)
		} else {
			f.parseSection(subsection, maxRecursion-1)
		}
	}

	for _, field := range schema.Fields {
		if _, ok := f.elements[field.ID]; ok {
			err := common.NewFieldValidatorError(field.ID)
			err.AddFailedConstraint(common.FailedConstraintError{
				Constraint: "unique id",
				Message:    "id must be unique",
			})
			f.errors = append(f.errors, err)
		}

		f.elements[field.ID] = schemaElement{FIELD, field}

		for _, subfield := range field.Subfields {
			if _, ok := f.elements[subfield.ID]; ok {
				err := common.NewFieldValidatorError(subfield.ID)
				err.AddFailedConstraint(common.FailedConstraintError{
					Constraint: "unique id",
					Message:    "id must be unique",
				})
				f.errors = append(f.errors, err)
			}
			f.elements[subfield.ID] = schemaElement{SUBFIELD, subfield}
		}
	}
}

func (f *FormValidationHelper) getElement(id int, elementType schemaElementType) (schemaElement, error) {
	if element, ok := f.elements[id]; ok {
		if element.Type != elementType {
			return schemaElement{}, ErrWrongElementType
		}

		return element, nil
	}
	return schemaElement{}, ErrElementNotFound
}

func (f *FormValidationHelper) GetSection(id int) (*model.WebFormSection, error) {
	element, err := f.getElement(id, SECTION)

	if err != nil {
		return nil, err
	}

	return element.Element.(*model.WebFormSection), nil
}

func (f *FormValidationHelper) GetField(id int) (*model.WebFormField, error) {
	element, err := f.getElement(id, FIELD)

	if err != nil {
		return nil, err
	}

	return element.Element.(*model.WebFormField), nil
}

func (f *FormValidationHelper) GetSubfield(id int) (*model.WebFormSubfield, error) {
	element, err := f.getElement(id, SUBFIELD)

	if err != nil {
		// check if subfield exists as field
		if element, err := f.getElement(id, FIELD); err == nil {
			return element.Element.(*model.WebFormField).WebFormSubfield, nil
		}

		return nil, err
	}

	return element.Element.(*model.WebFormSubfield), nil
}
