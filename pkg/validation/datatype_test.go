package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedDatatypes = []Datatype{
	{
		definition: DatatypeDefinition{Identifier: "@test",
			DisplayName:     "test",
			AllowsSubfields: false},
	},
	{
		definition: DatatypeDefinition{
			Identifier:      "@string",
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

	// try adding a datatype with InheritsFrom
	inheritanceDatatype := Datatype{
		definition: DatatypeDefinition{
			Identifier:      "inheritingDatatype",
			DisplayName:     "test",
			AllowsSubfields: false,
			InheritsFrom:    "@nonexistent",
		},
	}

	err = repository.AddDatatype(&inheritanceDatatype)
	assert.Error(t, err)
	assert.Equal(t, ErrDatatypeInvalidParent, err)

	inheritanceDatatype.definition.InheritsFrom = expectedDatatypes[0].definition.Identifier
	assert.NoError(t, repository.AddDatatype(&inheritanceDatatype))
}

func TestDatatypeRepository_DeleteDatatype(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	for _, datatypeDefinition := range expectedDatatypes {
		err := repository.DeleteDatatype(datatypeDefinition.definition.Identifier)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.DeleteDatatype(datatypeDefinition.definition.Identifier))
	}

	for _, datatypeDefinition := range expectedDatatypes {
		err := repository.DeleteDatatype(datatypeDefinition.definition.Identifier)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}

	// try removing a datatype where a different datatype inherits from
	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	inheritsFrom := "@string"
	inheritanceDatatype := Datatype{
		definition: DatatypeDefinition{
			Identifier:      "inheritingDatatype",
			DisplayName:     "test",
			AllowsSubfields: false,
			InheritsFrom:    inheritsFrom,
		},
	}
	assert.NoError(t, repository.AddDatatype(&inheritanceDatatype))

	err = repository.DeleteDatatype(inheritsFrom)
	assert.Error(t, err)
	assert.Equal(t, ErrDatatypeIsParent, err)

	assert.NoError(t, repository.DeleteDatatype(inheritanceDatatype.definition.Identifier))
	assert.NoError(t, repository.DeleteDatatype(inheritsFrom))
}

func TestDatatypeRepository_GetDatatype(t *testing.T) {
	repository, err := NewDatatypeRepository()
	assert.Nil(t, err)

	for _, expectedDatatype := range expectedDatatypes {
		datatype, err := repository.GetDatatype(expectedDatatype.definition.Identifier)
		assert.Nil(t, datatype)
		assert.Error(t, err)
		assert.Equal(t, ErrDatatypeNotFound, err)
	}

	for _, datatypeDefinition := range expectedDatatypes {
		assert.NoError(t, repository.AddDatatype(&datatypeDefinition))
	}

	for _, expectedDatatype := range expectedDatatypes {
		datatype, err := repository.GetDatatype(expectedDatatype.definition.Identifier)
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
		expected[expectedDatatype.definition.Identifier] = expectedDatatype
	}

	for _, definition := range result {
		res, ok := expected[definition.Identifier]
		assert.True(t, ok)
		assert.NotEmpty(t, res)
		assert.Equal(t, res.definition, definition)
	}
}
