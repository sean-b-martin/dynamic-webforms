package model

import "encoding/json"

// WebFormSchema contains all Sections of a form. ID is not used directly and can get set freely.
type WebFormSchema struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Sections []*WebFormSection `json:"sections"`
}

// WebFormSection contains 0 to n subsections and 0 to n fields. ResourceLinks are only used in the frontend to load
// resources like files and images.
type WebFormSection struct {
	ID            int               `json:"id"`
	Title         string            `json:"name"`
	Description   []string          `json:"description"`
	ResourceLinks []string          `json:"resourceLinks"`
	Subsections   []*WebFormSection `json:"subsections"`
	Fields        []*WebFormField   `json:"fields"`
}

// WebFormField is a WebFormSubfield with optional 0 to n Subfields. An example for Subfields are table columns.
// Normal fields like textboxes do not allow the use of Subfields.
type WebFormField struct {
	*WebFormSubfield
	Subfields []*WebFormSubfield `json:"subfields"`
}

// WebFormSubfield contains all metadata for an input field in a web form.
type WebFormSubfield struct {
	ID               int                      `json:"id"`
	Title            string                   `json:"name"`
	Type             string                   `json:"type"`
	Description      []string                 `json:"description"`
	ValidationSchema *WebFormValidationSchema `json:"validationSchema"`
}

// WebFormValidationSchema contains validation rules that are supported by all datatypes and dynamic validation rules
// that are only available to some datatypes.
type WebFormValidationSchema struct {
	BasicConstraints
	DynamicConstraints json.RawMessage `json:"dynamicConstraints"`
}

// BasicConstraints contains all attributes for validation that are used for any datatype.
type BasicConstraints struct {
	MinItems *int `json:"minItems"`
	MaxItems *int `json:"maxItems"`
}

// GenerateIDs generates all ID's used for validating user input against the form schema and validation rules.
// The function sets all ID's of the contained Sections, subsections and fields.
func (w *WebFormSchema) GenerateIDs() {
	currentNumber := 1
	for _, section := range w.Sections {
		currentNumber = section.GenerateIDs(currentNumber)
	}
}

// GenerateIDs Generates ID's for the section itself, the Subsections and Fields.
// ID's should only be used for validation of user input.
func (w *WebFormSection) GenerateIDs(startingNumber int) int {
	currentNumber := startingNumber

	w.ID = currentNumber
	currentNumber++

	for _, section := range w.Subsections {
		currentNumber = section.GenerateIDs(currentNumber)
	}

	for _, field := range w.Fields {
		currentNumber = field.GenerateIDs(currentNumber)
	}

	return currentNumber
}

// GenerateIDs Generates ID's for WebFormField and Subfields.
func (w *WebFormField) GenerateIDs(startingNumber int) int {
	currentNumber := startingNumber
	w.ID = currentNumber
	currentNumber++

	for _, field := range w.Subfields {
		currentNumber = field.GenerateIDs(currentNumber)
	}

	return currentNumber
}

// GenerateIDs Generates the ID for the current field.
func (w *WebFormSubfield) GenerateIDs(startingNumber int) int {
	w.ID = startingNumber
	startingNumber++
	return startingNumber
}
