package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWebForm_GenerateIDs(t *testing.T) {
	form := WebFormSchema{ID: "test", Sections: make([]*WebFormSection, 0), Name: "Test"}
	form.Sections = append(form.Sections, &WebFormSection{Title: "Section1"}, &WebFormSection{Title: "Section2"})
	form.Sections[1].Subsections = append(form.Sections[0].Subsections, &WebFormSection{Title: "Subsection1"},
		&WebFormSection{Title: "Subsection2"})

	form.Sections[1].Subsections[0].Fields = append(form.Sections[1].Subsections[0].Fields,
		&WebFormFieldContainer{WebFormField: &WebFormField{Title: "Field1"}},
		&WebFormFieldContainer{WebFormField: &WebFormField{Title: "Field2"}},
		&WebFormFieldContainer{WebFormField: &WebFormField{Title: "Field3"}},
		&WebFormFieldContainer{WebFormField: &WebFormField{Title: "Field4"}, Subfields: make([]*WebFormField, 0)})

	form.Sections[1].Subsections[0].Fields[3].Subfields = append(form.Sections[1].Subsections[0].Fields[3].Subfields,
		&WebFormField{Title: "Column1"}, &WebFormField{Title: "Column2"})

	form.GenerateIDs()
	currentValue := 1
	sections := form.Sections

	assert.Equal(t, sections[0].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].ID, currentValue)
	currentValue++

	// evaluate subsection 1
	assert.Equal(t, sections[1].Subsections[0].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].Subsections[0].Fields[0].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].Subsections[0].Fields[1].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].Subsections[0].Fields[2].ID, currentValue)
	currentValue++

	// evaluate complex type
	assert.Equal(t, sections[1].Subsections[0].Fields[3].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].Subsections[0].Fields[3].Subfields[0].ID, currentValue)
	currentValue++
	assert.Equal(t, sections[1].Subsections[0].Fields[3].Subfields[1].ID, currentValue)
	currentValue++

	assert.Equal(t, sections[1].Subsections[1].ID, currentValue)
}

func TestWebFormSection_GenerateIDs(t *testing.T) {
	section := WebFormSection{Subsections: make([]*WebFormSection, 0)}
	section.Subsections = append(section.Subsections, &WebFormSection{Title: "Subsection1"},
		&WebFormSection{Title: "Subsection2"})
	section.GenerateIDs(0)

	assert.Equal(t, 0, section.ID)
	assert.Equal(t, 1, section.Subsections[0].ID)
	assert.Equal(t, 2, section.Subsections[1].ID)
}

func TestWebFormFieldContainer_GenerateIDs(t *testing.T) {
	field := WebFormFieldContainer{WebFormField: &WebFormField{}, Subfields: make([]*WebFormField, 0)}
	field.Subfields = append(field.Subfields, &WebFormField{Title: "Column1"}, &WebFormField{Title: "Column2"})

	field.GenerateIDs(0)
	assert.Equal(t, 0, field.ID)
	assert.Equal(t, 1, field.Subfields[0].ID)
	assert.Equal(t, 2, field.Subfields[1].ID)
}

func TestWebFormField_GenerateIDs(t *testing.T) {
	field := WebFormField{}
	field.GenerateIDs(0)
	assert.Equal(t, 0, field.ID)

	field.GenerateIDs(1)
	assert.Equal(t, 1, field.ID)
}
