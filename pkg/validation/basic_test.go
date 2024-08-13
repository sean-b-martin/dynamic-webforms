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
	validator := NewBasicConstraintsValidator(1)
	basicConstraints := model.BasicConstraints{
		MinItems: initPointer(0),
		MaxItems: initPointer(2),
	}

	data := &model.WebFormDataRaw{
		SchemaElementID: 0,
		Data:            nil,
	}

	assert.Empty(t, validator.ValidateData(data, &basicConstraints))

	data.Data = []json.RawMessage{json.RawMessage(`1`), json.RawMessage(`2`), json.RawMessage(`3`)}
	assert.NotEmpty(t, validator.ValidateData(data, &basicConstraints))
}

func TestBasicConstraintsValidator_ValidateSchema(t *testing.T) {
	validator := NewBasicConstraintsValidator(1)

	field := &model.WebFormField{
		WebFormSubfield: &model.WebFormSubfield{
			ID:          0,
			Title:       "",
			Type:        "",
			Description: nil,
			ValidationSchema: &model.WebFormValidationSchema{
				BasicConstraints: model.BasicConstraints{},
			},
		},
		Subfields: []*model.WebFormSubfield{
			{
				ID:          0,
				Title:       "",
				Type:        "",
				Description: nil,
				ValidationSchema: &model.WebFormValidationSchema{
					BasicConstraints: model.BasicConstraints{},
				},
			},
		},
	}

	assert.Empty(t, validator.ValidateSchema(field))

	field.ValidationSchema.BasicConstraints.MaxItems = initPointer(-1)
	assert.NotEmpty(t, validator.ValidateSchema(field))

	// MinItems > MaxItems
	field.ValidationSchema.BasicConstraints.MaxItems = initPointer(2)
	field.ValidationSchema.BasicConstraints.MinItems = initPointer(3)
	assert.NotEmpty(t, validator.ValidateSchema(field))

	// MinItems == MaxItems
	field.ValidationSchema.BasicConstraints.MinItems = initPointer(2)
	assert.Empty(t, validator.ValidateSchema(field))

	// Error in Subfield
	field.Subfields[0].ValidationSchema.BasicConstraints.MaxItems = initPointer(-1)
	assert.NotEmpty(t, validator.ValidateSchema(field))
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
