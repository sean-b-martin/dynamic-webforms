package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWebForm_GenerateIDs(t *testing.T) {
	form := WebFormSchema{Sections: make([]*WebFormSection, 0), Title: "Test"}
	form.Sections = append(form.Sections, &WebFormSection{Title: "Section1"}, &WebFormSection{Title: "Section2"})
	form.Sections[1].Subsections = append(form.Sections[0].Subsections, &WebFormSection{Title: "Subsection1"},
		&WebFormSection{Title: "Subsection2"})

	form.Sections[1].Subsections[0].Fields = append(form.Sections[1].Subsections[0].Fields,
		&WebFormField{WebFormSubfield: &WebFormSubfield{Title: "Field1"}},
		&WebFormField{WebFormSubfield: &WebFormSubfield{Title: "Field2"}},
		&WebFormField{WebFormSubfield: &WebFormSubfield{Title: "Field3"}},
		&WebFormField{WebFormSubfield: &WebFormSubfield{Title: "Field4"}, Subfields: make([]*WebFormSubfield, 0)})

	form.Sections[1].Subsections[0].Fields[3].Subfields = append(form.Sections[1].Subsections[0].Fields[3].Subfields,
		&WebFormSubfield{Title: "Column1"}, &WebFormSubfield{Title: "Column2"})

	form.GenerateIDs()
	currentValue := 1
	sections := form.Sections

	assert.Equal(t, currentValue, sections[0].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].ID)
	currentValue++

	// evaluate subsection 1
	assert.Equal(t, currentValue, sections[1].Subsections[0].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[0].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[1].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[2].ID)
	currentValue++

	// evaluate complex type
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[3].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[3].Subfields[0].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[0].Fields[3].Subfields[1].ID)
	currentValue++
	assert.Equal(t, currentValue, sections[1].Subsections[1].ID)
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
	field := WebFormField{WebFormSubfield: &WebFormSubfield{}, Subfields: make([]*WebFormSubfield, 0)}
	field.Subfields = append(field.Subfields, &WebFormSubfield{Title: "Column1"}, &WebFormSubfield{Title: "Column2"})

	field.GenerateIDs(0)
	assert.Equal(t, 0, field.ID)
	assert.Equal(t, 1, field.Subfields[0].ID)
	assert.Equal(t, 2, field.Subfields[1].ID)
}

func TestWebFormField_GenerateIDs(t *testing.T) {
	field := WebFormSubfield{}
	field.GenerateIDs(0)
	assert.Equal(t, 0, field.ID)

	field.GenerateIDs(1)
	assert.Equal(t, 1, field.ID)
}
