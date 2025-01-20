package validator

import (
	"encoding/json"
)

func unmarshalValidationSchema[T any](data *json.RawMessage, result *T) *FailedConstraintError {
	if err := json.Unmarshal(*data, result); err != nil {
		return &FailedConstraintError{
			Constraint: "schema",
			DataIndex:  -1,
			Message:    err.Error(),
			Config:     nil,
		}
	}

	return nil
}
