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
	ID            int                      `json:"id"`
	Title         string                   `json:"name"`
	Description   []string                 `json:"description"`
	ResourceLinks []string                 `json:"resourceLinks"`
	Subsections   []*WebFormSection        `json:"subsections"`
	Fields        []*WebFormFieldContainer `json:"fields"`
}

// WebFormFieldContainer contains a WebFormField and 0 to n Subfields. An example for Subfields are table columns.
// Normal fields like textboxes do not allow the use of Subfields.
type WebFormFieldContainer struct {
	*WebFormField
	Subfields []*WebFormField `json:"subfields"`
}

// WebFormField contains all metadata for an input field in a web form.
type WebFormField struct {
	ID               int                     `json:"id"`
	Title            string                  `json:"name"`
	Type             string                  `json:"type"`
	Description      []string                `json:"description"`
	ValidationSchema WebFormValidationSchema `json:"validationSchema"`
}

// WebFormValidationSchema contains validation rules that are supported by all datatypes and dynamic validation rules
// that are only available to some datatypes.
type WebFormValidationSchema struct {
	BasicConstraints
	DynamicConstraints map[string]json.RawMessage `json:"dynamicConstraints"`
}

// BasicConstraints contains all attributes for validation that are used for any datatype.
type BasicConstraints struct {
	MinItems *int `json:"minItems"`
	MaxItems *int `json:"maxItems"`
}
