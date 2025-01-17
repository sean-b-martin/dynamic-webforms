package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedDatatypes = []Datatype{
	{
		definition: DatatypeDefinition{ID: "@test",
			DisplayName:     "test",
			AllowsSubfields: false},
	},
	{
		definition: DatatypeDefinition{
			ID:              "@string",
			DisplayName:     "string",
			AllowsSubfields: false,
		},
	},
}

func TestDatatypeRepository_AddDatatype(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, datatypeDefinition := range expectedDatatypes {
		err := repository.AddDatatype(&datatypeDefinition)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeDuplicate, err)
	}
}

func TestDatatypeRepository_DeleteDatatype(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	for _, datatypeDefinition := range expectedDatatypes {
		err := repository.DeleteDatatype(datatypeDefinition.definition.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.DeleteDatatype(datatypeDefinition.definition.ID))
	}

	for _, datatypeDefinition := range expectedDatatypes {
		err := repository.DeleteDatatype(datatypeDefinition.definition.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}
}

func TestDatatypeRepository_GetDatatype(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	for _, expectedDatatype := range expectedDatatypes {
		datatype, err := repository.GetDatatype(expectedDatatype.definition.ID)
		assert.Nil(t, datatype)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, expectedDatatype := range expectedDatatypes {
		datatype, err := repository.GetDatatype(expectedDatatype.definition.ID)
		if assert.NotNil(t, datatype) {
			assert.Equal(t, expectedDatatype, *datatype)
		}
		assert.NoError(t, err)
	}
}

func TestDatatypeRepository_GetDatatypeDefinitions(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	result := repository.GetDatatypeDefinitions()
	assert.Equal(t, 0, len(result))

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	result = repository.GetDatatypeDefinitions()
	assert.Equal(t, len(expectedDatatypes), len(result))

	expected := make(map[string]Datatype)

	for _, expectedDatatype := range expectedDatatypes {
		expected[expectedDatatype.definition.ID] = expectedDatatype
	}

	for _, definition := range result {
		res, ok := expected[definition.ID]
		assert.True(t, ok)
		assert.NotEmpty(t, res)
		assert.Equal(t, res.definition, definition)
	}
}
