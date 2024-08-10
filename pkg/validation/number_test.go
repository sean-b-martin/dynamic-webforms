package validation

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_numberTypeValidator_ValidateData(t *testing.T) {
	validator := numberTypeValidator[int]{}
	dynamicConstraints := DynamicNumberValidationSchema[int]{
		Lt: initPointer(100),
		Gt: initPointer(1),
	}

	j, _ := json.Marshal(dynamicConstraints)
	validationSchema := &model.WebFormValidationSchema{
		BasicConstraints:   model.BasicConstraints{},
		DynamicConstraints: j,
	}

	var data *model.WebFormDataRaw
	data = &model.WebFormDataRaw{
		SchemaElementID: 0,
		Data:            []json.RawMessage{json.RawMessage(`2`)},
	}

	// no error
	assert.Empty(t, validator.ValidateData(data, validationSchema))

	// data == gt
	data.Data[0] = json.RawMessage(`1`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// data < gt
	data.Data[0] = json.RawMessage(`-1`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// data == lt
	data.Data[0] = json.RawMessage(`100`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// data > lt
	data.Data[0] = json.RawMessage(`101`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// gte & lte

	validationSchema.DynamicConstraints, _ = json.Marshal(&DynamicNumberValidationSchema[int]{
		Lte: initPointer(5),
		Gte: initPointer(1),
	})

	// data == gte
	data.Data[0] = json.RawMessage(`1`)
	assert.Empty(t, validator.ValidateData(data, validationSchema))

	// data < gte
	data.Data[0] = json.RawMessage(`-1`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// data == lte
	data.Data[0] = json.RawMessage(`5`)
	assert.Empty(t, validator.ValidateData(data, validationSchema))

	// data >= lte
	data.Data[0] = json.RawMessage(`100`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	validationSchema.DynamicConstraints, _ = json.Marshal(&DynamicNumberValidationSchema[int]{
		MinDigits: initPointer(2),
		MaxDigits: initPointer(3),
	})

	// len(data) < MinDigits
	data.Data[0] = json.RawMessage(`1`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))

	// len(data) == MinDigits
	data.Data[0] = json.RawMessage(`10`)
	assert.Empty(t, validator.ValidateData(data, validationSchema))

	// len(data) == MaxDigits
	data.Data[0] = json.RawMessage(`100`)
	assert.Empty(t, validator.ValidateData(data, validationSchema))

	// len(data) > MaxDigits
	data.Data[0] = json.RawMessage(`1000`)
	assert.NotEmpty(t, validator.ValidateData(data, validationSchema))
}

func Test_numberTypeValidator_ValidateSchema(t *testing.T) {
	validator := numberTypeValidator[int]{}

	field := &model.WebFormField{
		WebFormSubfield: &model.WebFormSubfield{
			ID:               0,
			Title:            "",
			Type:             "",
			Description:      nil,
			ValidationSchema: &model.WebFormValidationSchema{},
		},
		Subfields: nil,
	}

	assert.Empty(t, validator.ValidateSchema(field))

	dynamicConstraints := make(map[string]interface{})

	field.ValidationSchema.DynamicConstraints, _ = json.Marshal(DynamicNumberValidationSchema[int]{
		Lt:        initPointer(1),
		Gt:        initPointer(1),
		Lte:       initPointer(1),
		Gte:       initPointer(1),
		MinDigits: initPointer(1),
		MaxDigits: initPointer(1),
	})

	assert.Empty(t, validator.ValidateSchema(field))

	// unknown keys should not lead to errors
	dynamicConstraints["newKey"] = 2
	dynamicConstraints["newKey2"] = 2.5
	field.ValidationSchema.DynamicConstraints, _ = json.Marshal(dynamicConstraints)
	assert.Empty(t, validator.ValidateSchema(field))

	// wrong type
	dynamicConstraints["lt"] = "hello"
	field.ValidationSchema.DynamicConstraints, _ = json.Marshal(dynamicConstraints)
	assert.NotEmpty(t, validator.ValidateSchema(field))
}

func Test_bigNumberTypeValidator_ValidateData(t *testing.T) {

}

func Test_bigNumberTypeValidator_ValidateSchema(t *testing.T) {

}

func initPointer[T Numbers](value T) *T {
	return &value
}
