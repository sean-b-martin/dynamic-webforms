package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var datatypeDefinitions = []DatatypeDefinition{
	{
		Identifier:      "@test",
		DisplayName:     "test",
		AllowsSubfields: false,
	},
	{
		Identifier:      "@string",
		DisplayName:     "string",
		AllowsSubfields: false,
	},
}

func TestDatatypeRepository_AddDatatype(t *testing.T) {
	repository := NewDatatypeRepository()

	for _, datatypeDefinition := range datatypeDefinitions {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		err := repository.AddDatatype(&datatypeDefinition)
		assert.Error(t, err)
		assert.Equal(t, DatatypeDuplicateError, err)
	}
}

func TestDatatypeRepository_DeleteDatatype(t *testing.T) {
	repository := NewDatatypeRepository()

	for _, datatypeDefinition := range datatypeDefinitions {
		err := repository.DeleteDatatype(datatypeDefinition.Identifier)
		assert.Error(t, err)
		assert.Equal(t, DatatypeNotFoundError, err)
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		err := repository.AddDatatype(&datatypeDefinition)
		assert.NoError(t, err)
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		assert.NoError(t, repository.DeleteDatatype(datatypeDefinition.Identifier))
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		err := repository.DeleteDatatype(datatypeDefinition.Identifier)
		assert.Error(t, err)
		assert.Equal(t, DatatypeNotFoundError, err)
	}
}

func TestDatatypeRepository_GetDatatype(t *testing.T) {
	repository := NewDatatypeRepository()

	for _, datatypeDefinition := range datatypeDefinitions {
		datatype, err := repository.GetDatatype(datatypeDefinition.Identifier)
		assert.Nil(t, datatype)
		assert.Error(t, err)
		assert.Equal(t, DatatypeNotFoundError, err)
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, datatypeDefinition := range datatypeDefinitions {
		datatype, err := repository.GetDatatype(datatypeDefinition.Identifier)
		if assert.NotNil(t, datatype) {
			assert.Equal(t, datatypeDefinition, datatype.definition)
		}
		assert.NoError(t, err)
	}
}

func TestDatatypeRepository_GetDatatypeDefinitions(t *testing.T) {
	repository := NewDatatypeRepository()

	result := repository.GetDatatypeDefinitions()
	assert.Equal(t, 0, len(result))

	for _, datatypeDefinition := range datatypeDefinitions {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	result = repository.GetDatatypeDefinitions()
	assert.Equal(t, len(datatypeDefinitions), len(result))

	expected := make(map[string]DatatypeDefinition)

	for _, definition := range datatypeDefinitions {
		expected[definition.Identifier] = definition
	}

	for _, definition := range result {
		res, ok := expected[definition.Identifier]
		assert.True(t, ok)
		assert.NotEmpty(t, res)
		assert.Equal(t, res, definition)
	}
}
