package common

import (
	"encoding/json"
)

func UnmarshalValidationSchema[T any](data *json.RawMessage, result *T) *FailedConstraintError {
	if err := json.Unmarshal(*data, result); err != nil {
		return &FailedConstraintError{
			Constraint: "schema",
			DataIndex:  nil,
			Message:    err.Error(),
			Config:     nil,
		}
	}

	return nil
}
