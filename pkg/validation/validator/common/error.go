package common

import (
	"strconv"
	"strings"
)

type FailedConstraintError struct {
	Constraint string      `json:"constraint"`
	DataIndex  *int        `json:"dataIndex,omitempty"`
	Message    string      `json:"message"`
	Config     interface{} `json:"config,omitempty"`
}

func (e *FailedConstraintError) Error() string {
	return e.Message
}

type ValidatorError struct {
	SchemaElementID   int                     `json:"schemaElementID"`
	FailedConstraints []FailedConstraintError `json:"errors"`
}

func (e *ValidatorError) AddFailedConstraint(constraintError FailedConstraintError) {
	e.FailedConstraints = append(e.FailedConstraints, constraintError)
}

func (e *ValidatorError) IsEmpty() bool {
	return len(e.FailedConstraints) == 0
}

func NewFieldValidatorError(schemaElementID int) ValidatorError {
	return ValidatorError{SchemaElementID: schemaElementID}
}

func (e *ValidatorError) Error() string {
	if len(e.FailedConstraints) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for i, constraint := range e.FailedConstraints {
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(":")
		sb.WriteString(constraint.Error())
		sb.WriteString(",")
	}
	return sb.String()
}
