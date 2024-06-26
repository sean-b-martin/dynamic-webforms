package validation

import (
	"errors"
	"strconv"
	"strings"
)

var (
	DatatypeDuplicateError = errors.New("dynamic-forms: datatype already exists")
	DatatypeNotFoundError  = errors.New("dynamic-forms: datatype not found")
)

type FormValidationError struct {
	Name            string `json:"name"`
	SchemaElementID int    `json:"schemaElementID"`
	Message         string `json:"message"`
}

func (f FormValidationError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("form-validation: ")
	sb.WriteString(f.Name)
	sb.WriteString(strconv.Itoa(f.SchemaElementID))
	sb.WriteString(f.Message)
	return sb.String()
}
