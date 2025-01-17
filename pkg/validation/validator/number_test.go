package validator

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericNumberValidator_ValidateData(t *testing.T) {
	validator := &GenericNumberValidator[int]{}

	data := model.WebFormDataRaw{
		SchemaElementID: 1,
		Data:            intoRawMessage(128),
	}

	constraints := DynamicNumberValidationSchema[int]{
		Lt:        makePtr(1024),
		Gt:        nil,
		Lte:       nil,
		Gte:       nil,
		MinDigits: nil,
		MaxDigits: nil,
	}

	// value < lt
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value > lt
	constraints.Lt = makePtr(64)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value == lt
	constraints.Lt = makePtr(128)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)
	constraints.Lt = nil

	// value > gt
	constraints.Gt = makePtr(0)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value < gt
	constraints.Gt = makePtr(256)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value == gt
	constraints.Gt = makePtr(128)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)
	constraints.Gt = nil

	// value < lte
	constraints.Lte = makePtr(256)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value > lte
	constraints.Lte = makePtr(64)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value == lte
	constraints.Lte = makePtr(128)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)
	constraints.Lte = nil

	// value < gte
	constraints.Gte = makePtr(256)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value > gte
	constraints.Gte = makePtr(0)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value == gte
	constraints.Gte = makePtr(128)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)
	constraints.Gte = nil

	// value < MinDigits
	constraints.MinDigits = makePtr(5)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value > MinDigits
	constraints.MinDigits = makePtr(2)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value == MinDigits
	constraints.MinDigits = makePtr(3)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)
	constraints.MinDigits = nil

	// value < MaxDigits
	constraints.MaxDigits = makePtr(5)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)

	// value > MaxDigits
	constraints.MaxDigits = makePtr(2)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], true)

	// value == MaxDigits
	constraints.MaxDigits = makePtr(3)
	testValidate[int](t, validator, &data, intoRawMessage(constraints)[0], false)
}

func testValidate[T Number](t *testing.T, validator *GenericNumberValidator[T], data *model.WebFormDataRaw, constraints json.RawMessage, expectErr bool) {
	err := validator.ValidateData(data, intoRawMessage(constraints)[0])
	if expectErr {
		assert.False(t, err.IsEmpty())
		assert.Len(t, err.FailedConstraints, 1)
	} else {
		assert.True(t, err.IsEmpty())
	}
}

func intoRawMessage(data ...interface{}) []json.RawMessage {
	result := make([]json.RawMessage, 0, len(data))
	for _, d := range data {
		j, _ := json.Marshal(d)
		result = append(result, j)
	}

	return result
}

func makePtr[T any](value T) *T {
	return &value
}
