package validator

import (
	"fmt"
	"strings"
)

type FieldError struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
	Msg   string `json:"msg"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field '%s': %d", e.Name, e.Index)
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Empty() bool {
	return len(e.Errors) == 0
}

func (e *ValidationError) AddFieldError(f FieldError) {
	e.Errors = append(e.Errors, f)
}

func (e *ValidationError) Error() string {
	sb := strings.Builder{}
	for _, f := range e.Errors {
		sb.WriteString(f.Error())
	}
	return sb.String()
}
