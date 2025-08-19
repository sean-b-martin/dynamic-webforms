package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm_GenerateIDs(t *testing.T) {
	form := FormSchema{Sections: make([]*Section, 0), Title: "Test"}
	form.Sections = append(form.Sections, &Section{Title: "Section1"}, &Section{Title: "Section2"})
	form.Sections[1].Subsections = append(form.Sections[0].Subsections, &Section{Title: "Subsection1"},
		&Section{Title: "Subsection2"})

	form.Sections[1].Subsections[0].Fields = append(form.Sections[1].Subsections[0].Fields,
		&Field{Subfield: &Subfield{Title: "Field1"}},
		&Field{Subfield: &Subfield{Title: "Field2"}},
		&Field{Subfield: &Subfield{Title: "Field3"}},
		&Field{Subfield: &Subfield{Title: "Field4"}, Subfields: make([]*Subfield, 0)})

	form.Sections[1].Subsections[0].Fields[3].Subfields = append(form.Sections[1].Subsections[0].Fields[3].Subfields,
		&Subfield{Title: "Column1"}, &Subfield{Title: "Column2"})

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

func TestSection_GenerateIDs(t *testing.T) {
	section := Section{Subsections: make([]*Section, 0)}
	section.Subsections = append(section.Subsections, &Section{Title: "Subsection1"},
		&Section{Title: "Subsection2"})
	section.generateIDs(0)

	assert.Equal(t, 0, section.ID)
	assert.Equal(t, 1, section.Subsections[0].ID)
	assert.Equal(t, 2, section.Subsections[1].ID)
}

func TestField_GenerateIDs(t *testing.T) {
	field := Field{Subfield: &Subfield{}, Subfields: make([]*Subfield, 0)}
	field.Subfields = append(field.Subfields, &Subfield{Title: "Column1"}, &Subfield{Title: "Column2"})

	field.generateIDs(0)
	assert.Equal(t, 0, field.ID)
	assert.Equal(t, 1, field.Subfields[0].ID)
	assert.Equal(t, 2, field.Subfields[1].ID)
}

func TestSubfield_GenerateIDs(t *testing.T) {
	field := Subfield{}
	field.generateIDs(0)
	assert.Equal(t, 0, field.ID)

	field.generateIDs(1)
	assert.Equal(t, 1, field.ID)
}
