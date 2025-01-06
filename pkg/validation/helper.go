package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
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
	elements map[int]schemaElement
	errors   []FormValidationError
}

func NewFormValidationHelper() *FormValidationHelper {
	return &FormValidationHelper{
		elements: make(map[int]schemaElement),
		errors:   nil,
	}
}

func (f *FormValidationHelper) ParseForm(schema *model.WebFormSchema, maxRecursion int) []FormValidationError {
	for _, section := range schema.Sections {
		f.parseSection(section, maxRecursion)
	}

	return f.errors
}

func (f *FormValidationHelper) parseSection(schema *model.WebFormSection, maxRecursion int) {
	if _, ok := f.elements[schema.ID]; ok {
		f.errors = append(f.errors, NewSchemaError(schema.ID, "duplicate element id"))
	}

	f.elements[schema.ID] = schemaElement{SECTION, schema}

	for _, subsection := range schema.Subsections {
		if maxRecursion == 0 {
			f.errors = append(f.errors, NewSchemaError(schema.ID, "too many subsection levels"))
		} else {
			f.parseSection(subsection, maxRecursion-1)
		}
	}

	for _, field := range schema.Fields {
		if _, ok := f.elements[field.ID]; ok {
			f.errors = append(f.errors, NewSchemaError(field.ID, "duplicate element id"))
		}

		f.elements[field.ID] = schemaElement{FIELD, field}

		for _, subfield := range field.Subfields {
			if _, ok := f.elements[subfield.ID]; ok {
				f.errors = append(f.errors, NewSchemaError(subfield.ID, "duplicate element id"))
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
