package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

var ErrTooShort = errors.New("value is shorter than the minimum length")

func LengthOfInt(value int) int {
	return int(math.Log10(float64(value))) + 1
}

type MinLengthValidator struct {
	MinLength int `json:"minLength"`
}

func (m *MinLengthValidator) NewValidator(attributes *json.RawMessage) (Validator, error) {
	if attributes == nil {
		return &MinLengthValidator{MinLength: 0}, nil
	}

	var validator MinLengthValidator
	if err := json.Unmarshal(*attributes, &validator); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	return &validator, nil
}

func (m *MinLengthValidator) Validate(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if utf8.RuneCountInString(v) < m.MinLength {
			return ErrTooShort
		}
	case int:
		if int(math.Log10(float64(v)))+1 < m.MinLength {
			return ErrTooShort
		}
	}

	return nil
}

func (m *MinLengthValidator) Info() Definition {
	return Definition{
		Description: "fails if the length is less then minLength",
		Attributes: map[string]any{
			"minLength": 0,
		},
	}
}
