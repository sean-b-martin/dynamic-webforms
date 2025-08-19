package model

import "encoding/json"

// FormSchema contains all Sections of a form
type FormSchema struct {
	Title    string     `json:"title"`
	Sections []*Section `json:"sections"`
}

// Section contains 0 to n subsections and 0 to n fields
type Section struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Subsections []*Section `json:"subsections"`
	Fields      []*Field   `json:"fields"`
}

// Field is a Subfield with optional 0 to n Subfields. An example for Subfields are table columns.
// Normal fields like textboxes do not allow the use of Subfields.
type Field struct {
	*Subfield
	Subfields []*Subfield `json:"subfields"`
}

// Subfield contains all metadata for an input field in a web form.
type Subfield struct {
	ID               int                   `json:"id"`
	Title            string                `json:"title"`
	Type             string                `json:"type"`
	Description      string                `json:"description"`
	ValidationSchema *FormValidationSchema `json:"validationSchema"`
}

// FormValidationSchema contains validation rules that are supported by all datatypes and dynamic validation rules
// that are only available to some datatypes.
type FormValidationSchema struct {
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
func (w *FormSchema) GenerateIDs() {
	currentNumber := 1
	for _, section := range w.Sections {
		currentNumber = section.generateIDs(currentNumber)
	}
}

// generateIDs Generates ID's for the section itself, the Subsections and Fields.
// ID's should only be used for validation of user input.
func (w *Section) generateIDs(startingNumber int) int {
	currentNumber := startingNumber

	w.ID = currentNumber
	currentNumber++

	for _, section := range w.Subsections {
		currentNumber = section.generateIDs(currentNumber)
	}

	for _, field := range w.Fields {
		currentNumber = field.generateIDs(currentNumber)
	}

	return currentNumber
}

// generateIDs Generates ID's for Field and Subfields.
func (w *Field) generateIDs(startingNumber int) int {
	currentNumber := startingNumber
	w.ID = currentNumber
	currentNumber++

	for _, field := range w.Subfields {
		currentNumber = field.generateIDs(currentNumber)
	}

	return currentNumber
}

// generateIDs Generates the ID for the current field.
func (w *Subfield) generateIDs(startingNumber int) int {
	w.ID = startingNumber
	startingNumber++
	return startingNumber
}
