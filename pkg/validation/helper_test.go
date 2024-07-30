package validation

import (
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var schema = model.WebFormSchema{
	ID:   "",
	Name: "",
	Sections: []*model.WebFormSection{
		{
			ID:            0,
			Title:         "section1",
			Description:   nil,
			ResourceLinks: nil,
			Subsections:   nil,
			Fields: []*model.WebFormField{
				{
					WebFormSubfield: &model.WebFormSubfield{
						ID:               2,
						Title:            "field1",
						Type:             "",
						Description:      nil,
						ValidationSchema: &model.WebFormValidationSchema{},
					},
					Subfields: []*model.WebFormSubfield{
						{
							ID:               3,
							Title:            "subfield1",
							Type:             "",
							Description:      nil,
							ValidationSchema: &model.WebFormValidationSchema{},
						},
						{
							ID:               4,
							Title:            "subfield2",
							Type:             "",
							Description:      nil,
							ValidationSchema: &model.WebFormValidationSchema{},
						},
					},
				},
			},
		},
		{
			ID:            1,
			Title:         "section2",
			Description:   nil,
			ResourceLinks: nil,
			Subsections: []*model.WebFormSection{{
				ID:            7,
				Title:         "subsection1",
				Description:   nil,
				ResourceLinks: nil,
				Subsections:   nil,
				Fields:        nil,
			}},
			Fields: []*model.WebFormField{
				{
					WebFormSubfield: &model.WebFormSubfield{
						ID:               5,
						Title:            "field2",
						Type:             "",
						Description:      nil,
						ValidationSchema: &model.WebFormValidationSchema{},
					},
					Subfields: nil,
				},
				{
					WebFormSubfield: &model.WebFormSubfield{
						ID:               6,
						Title:            "field3",
						Type:             "",
						Description:      nil,
						ValidationSchema: &model.WebFormValidationSchema{},
					},
					Subfields: nil,
				},
			},
		},
	},
}

func TestFormValidationHelper_GetField(t *testing.T) {
	helper := NewFormValidationHelper()
	helper.ParseForm(&schema, 5)

	field, err := helper.GetField(5)
	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, schema.Sections[1].Fields[0], field)

	// valid id, wrong type
	field, err = helper.GetField(1)
	assert.Error(t, err)
	assert.Nil(t, field)

	// invalid id
	field, err = helper.GetField(20)
	assert.Error(t, err)
	assert.Nil(t, field)
}

func TestFormValidationHelper_GetSection(t *testing.T) {
	helper := NewFormValidationHelper()
	helper.ParseForm(&schema, 5)

	field, err := helper.GetSection(1)
	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, schema.Sections[1], field)

	// valid id, wrong type
	field, err = helper.GetSection(5)
	assert.Error(t, err)
	assert.Nil(t, field)

	// invalid id
	field, err = helper.GetSection(20)
	assert.Error(t, err)
	assert.Nil(t, field)
}

func TestFormValidationHelper_GetSubfield(t *testing.T) {
	helper := NewFormValidationHelper()
	helper.ParseForm(&schema, 5)

	field, err := helper.GetSubfield(3)
	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, schema.Sections[0].Fields[0].Subfields[0], field)

	// id of field
	field, err = helper.GetSubfield(2)
	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, schema.Sections[0].Fields[0].WebFormSubfield, field)

	// valid id, wrong type
	field, err = helper.GetSubfield(1)
	assert.Error(t, err)
	assert.Nil(t, field)

	// invalid id
	field, err = helper.GetSubfield(20)
	assert.Error(t, err)
	assert.Nil(t, field)
}

func TestFormValidationHelper_ParseForm(t *testing.T) {
	helper := NewFormValidationHelper()
	assert.Empty(t, helper.elements)
	assert.Empty(t, helper.errors)

	helper.ParseForm(&schema, 5)
	assert.NotEmpty(t, helper.elements)
	assert.Empty(t, helper.errors)
	assert.Equal(t, schema.Sections[0], helper.elements[0].Element)
	assert.Equal(t, SECTION, helper.elements[0].Type)

	assert.Equal(t, schema.Sections[0].Fields[0], helper.elements[2].Element)
	assert.Equal(t, FIELD, helper.elements[2].Type)

	assert.Equal(t, schema.Sections[0].Fields[0].Subfields[0], helper.elements[3].Element)
	assert.Equal(t, SUBFIELD, helper.elements[3].Type)

	helper = NewFormValidationHelper()
	err := helper.ParseForm(&schema, 0)
	assert.NotEmpty(t, err)
	assert.Equal(t, "too many subsection levels", err[0].Message)

	// duplicate id's
	helper = NewFormValidationHelper()
	invalidSchema := model.WebFormSchema{
		ID:   "",
		Name: "",
		Sections: []*model.WebFormSection{
			{ID: 1, Title: "", Description: nil, ResourceLinks: nil, Subsections: nil, Fields: nil},
			{ID: 1, Title: "", Description: nil, ResourceLinks: nil, Subsections: nil, Fields: nil},
			{ID: 1, Title: "", Description: nil, ResourceLinks: nil, Subsections: nil, Fields: nil},
		},
	}

	err = helper.ParseForm(&invalidSchema, 5)
	assert.Len(t, err, 2)
}

func TestFormValidationHelper_getElement(t *testing.T) {
	helper := NewFormValidationHelper()
	helper.ParseForm(&schema, 5)

	element, err := helper.getElement(0, SECTION)
	assert.NoError(t, err)
	assert.Equal(t, SECTION, element.Type)
	assert.Equal(t, schema.Sections[0], element.Element)

	// wrong type
	element, err = helper.getElement(0, FIELD)
	assert.Error(t, err)
	assert.Nil(t, element.Element)

	// wrong id
	element, err = helper.getElement(20, SECTION)
	assert.Error(t, err)
	assert.Nil(t, element.Element)

	// field
	element, err = helper.getElement(2, FIELD)
	assert.NoError(t, err)
	assert.Equal(t, FIELD, element.Type)
	assert.Equal(t, schema.Sections[0].Fields[0], element.Element)

	// subfield
	element, err = helper.getElement(3, SUBFIELD)
	assert.NoError(t, err)
	assert.Equal(t, SUBFIELD, element.Type)
	assert.Equal(t, schema.Sections[0].Fields[0].Subfields[0], element.Element)
}

func TestFormValidationHelper_parseSection(t *testing.T) {
	section := &model.WebFormSection{
		ID:            0,
		Title:         "",
		Description:   nil,
		ResourceLinks: nil,
		Subsections:   nil,
		Fields: []*model.WebFormField{
			{
				WebFormSubfield: &model.WebFormSubfield{
					ID:               1,
					Title:            "",
					Type:             "",
					Description:      nil,
					ValidationSchema: &model.WebFormValidationSchema{},
				},
				Subfields: nil,
			},
		},
	}

	helper := NewFormValidationHelper()
	helper.parseSection(section, 5)

	assert.Equal(t, section, helper.elements[0].Element)
	assert.Equal(t, section.Fields[0], helper.elements[1].Element)

	section.Subsections = make([]*model.WebFormSection, 0)
	section.Subsections = append(section.Subsections, &model.WebFormSection{
		ID: 2,
		Subsections: []*model.WebFormSection{
			{
				ID: 3,
				Subsections: []*model.WebFormSection{
					{
						ID: 4, Subsections: []*model.WebFormSection{
							{ID: 5},
						},
					},
				},
			},
		},
	})

	helper = NewFormValidationHelper()
	helper.parseSection(section, 2)
	assert.NotEmpty(t, helper.errors)
	assert.Equal(t, "too many subsection levels", helper.errors[0].Message)

}

func TestNewFormValidationHelper(t *testing.T) {
	helper := NewFormValidationHelper()
	assert.NotNil(t, helper)
}
