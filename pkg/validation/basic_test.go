package validation

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBasicConstraintsValidator(t *testing.T) {
	assert.NotNil(t, NewBasicConstraintsValidator(1))
}

func TestBasicConstraintsValidator_ValidateData(t *testing.T) {

}

func TestBasicConstraintsValidator_ValidateSchema(t *testing.T) {

}

func TestBasicConstraintsValidator_validateItemCount(t *testing.T) {
	validator := NewBasicConstraintsValidator(1)
	basicConstraints := model.BasicConstraints{
		MinItems: initPointer(0),
		MaxItems: initPointer(3),
	}

	data := &model.WebFormDataRaw{
		SchemaElementID: 0,
		Data:            nil,
	}

	assert.Empty(t, validator.validateItemCount(data, &basicConstraints))

	// MaxItems == len(Data)
	data.Data = []json.RawMessage{json.RawMessage(`1`), json.RawMessage(`2`), json.RawMessage(`3`)}
	assert.Empty(t, validator.validateItemCount(data, &basicConstraints))

	// MaxItems < len(Data)
	data.Data = append(data.Data, json.RawMessage(`4`))
	assert.NotEmpty(t, validator.validateItemCount(data, &basicConstraints))

	// default value
	basicConstraints.MinItems = nil
	basicConstraints.MaxItems = nil
	assert.NotEmpty(t, validator.validateItemCount(data, &basicConstraints))

	validator = NewBasicConstraintsValidator(100)
	assert.Empty(t, validator.validateItemCount(data, &basicConstraints))
}
