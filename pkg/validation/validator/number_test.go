package validator

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericNumberValidator_ValidateFieldData(t *testing.T) {
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
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value > lt
	constraints.Lt = makePtr(64)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value == lt
	constraints.Lt = makePtr(128)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)
	constraints.Lt = nil

	// value > gt
	constraints.Gt = makePtr(0)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value < gt
	constraints.Gt = makePtr(256)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value == gt
	constraints.Gt = makePtr(128)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)
	constraints.Gt = nil

	// value < lte
	constraints.Lte = makePtr(256)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value > lte
	constraints.Lte = makePtr(64)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value == lte
	constraints.Lte = makePtr(128)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)
	constraints.Lte = nil

	// value < gte
	constraints.Gte = makePtr(256)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value > gte
	constraints.Gte = makePtr(0)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value == gte
	constraints.Gte = makePtr(128)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)
	constraints.Gte = nil

	// value < MinDigits
	constraints.MinDigits = makePtr(5)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value > MinDigits
	constraints.MinDigits = makePtr(2)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value == MinDigits
	constraints.MinDigits = makePtr(3)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)
	constraints.MinDigits = nil

	// value < MaxDigits
	constraints.MaxDigits = makePtr(5)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)

	// value > MaxDigits
	constraints.MaxDigits = makePtr(2)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], true)

	// value == MaxDigits
	constraints.MaxDigits = makePtr(3)
	testValidateData[int](t, &data, intoRawMessage(constraints)[0], false)
}

func Test_numberTypeValidator_Initialize(t *testing.T) {
	id := 0
	constraints := DynamicNumberValidationSchema[int]{
		Lt:        makePtr(2),
		Gt:        makePtr(1),
		Lte:       makePtr(2),
		Gte:       makePtr(1),
		MinDigits: makePtr(1),
		MaxDigits: makePtr(1),
	}

	testInitialize[int](t, id, intoRawMessage(constraints)[0], false)

	// unknown keys should not lead to errors
	m := make(map[string]interface{})
	m["invalid"] = 2
	testInitialize[int](t, id, intoRawMessage(m)[0], false)

	// wrong type
	m["lt"] = "hello"
	testInitialize[int](t, id, intoRawMessage(m)[0], true)

	// gt > lt
	constraints = DynamicNumberValidationSchema[int]{
		Lt: makePtr(2),
		Gt: makePtr(3),
	}
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)
	// gt == lt
	constraints.Gt = makePtr(3)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)
	constraints.Gt = nil

	// gte > lt
	constraints.Gte = makePtr(4)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)

	// gte == lt
	constraints.Gte = makePtr(3)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)
	constraints.Gte = nil
	constraints.Lt = nil

	// lte == gt
	constraints.Lte = makePtr(1)
	constraints.Gt = makePtr(1)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)

	// lte < gt
	constraints.Lte = makePtr(0)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)

	// lte < gte
	constraints.Gt = nil
	constraints.Lte = makePtr(1)
	constraints.Gte = makePtr(2)
	testInitialize[int](t, id, intoRawMessage(constraints)[0], true)
}

func testValidateData[T Number](t *testing.T, data *model.WebFormDataRaw, constraints json.RawMessage, expectErr bool) {
	validator := GenericNumberValidator[T]{}
	err := validator.Initialize(data.SchemaElementID, &constraints)
	assert.True(t, err.IsEmpty())
	err = validator.ValidateFieldData(data)
	if expectErr {
		assert.False(t, err.IsEmpty())
		assert.Len(t, err.FailedConstraints, 1)
	} else {
		assert.True(t, err.IsEmpty())
	}
}

func testInitialize[T Number](t *testing.T, id int, constraints json.RawMessage, expectErr bool) {
	validator := &GenericNumberValidator[T]{}
	err := validator.Initialize(id, &intoRawMessage(constraints)[0])
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
