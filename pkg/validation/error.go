package validation

import (
	"errors"
	"strconv"
	"strings"
)

var (
	DatatypeDuplicateError     = errors.New("dynamic-forms: datatype already exists")
	DatatypeNotFoundError      = errors.New("dynamic-forms: datatype not found")
	DatatypeInvalidParentError = errors.New("dynamic-forms: invalid parent datatype")
	DatatypeIsParentError      = errors.New("dynamic-forms: datatype is parent")
	ElementNotFoundError       = errors.New("dynamic-forms: element not found")
	WrongElementTypeError      = errors.New("dynamic-forms: wrong element type")
)

type FormValidationError struct {
	Name            string `json:"name"`
	SchemaElementID int    `json:"schemaElementID"`
	Message         string `json:"message"`
}

func NewSchemaError(schemaElementID int, message string) FormValidationError {
	return FormValidationError{Name: "form schema error", SchemaElementID: schemaElementID, Message: message}
}

func NewDataError(schemaElementID int, message string) FormValidationError {
	return FormValidationError{Name: "form data error", SchemaElementID: schemaElementID, Message: message}
}

func (f FormValidationError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("form-validation: ")
	sb.WriteString(f.Name)
	sb.WriteString(strconv.Itoa(f.SchemaElementID))
	sb.WriteString(f.Message)
	return sb.String()
}
