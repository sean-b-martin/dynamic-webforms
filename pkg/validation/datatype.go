package validation

// DatatypeDefinition contains metadata about a datatype such as the identifier for using the datatype
// and whether the datatype allows subfields or not.
type DatatypeDefinition struct {
	Identifier      string `json:"identifier"`
	DisplayName     string `json:"displayName"`
	AllowsSubfields bool   `json:"allowsSubfields"`
	InheritsFrom    string `json:"inheritsFrom"`
}

// Datatype contains its definition and all rules for validation
type Datatype struct {
	definition DatatypeDefinition
	DatatypeValidator
}

// DatatypeRepository contains the available datatypes.
type DatatypeRepository struct {
	datatypes map[string]*Datatype
}

func NewDatatype(definition DatatypeDefinition, validator DatatypeValidator) Datatype {
	return Datatype{definition: definition, DatatypeValidator: validator}
}

func NewDatatypeRepository() DatatypeRepository {
	return DatatypeRepository{datatypes: make(map[string]*Datatype)}
}

// GetDatatypeDefinitions returns a slice of all DatatypeDefinition available inside the repository
func (d *DatatypeRepository) GetDatatypeDefinitions() []DatatypeDefinition {
	result := make([]DatatypeDefinition, 0, len(d.datatypes))
	for _, datatype := range d.datatypes {
		result = append(result, datatype.definition)
	}

	return result
}

// GetDatatype returns the stored datatype,
// returns Datatype,nil if successful and nil, DatatypeNotFoundError when no datatype was found
func (d *DatatypeRepository) GetDatatype(identifier string) (*Datatype, error) {
	if datatype, ok := d.datatypes[identifier]; ok {
		return datatype, nil
	}

	return nil, DatatypeNotFoundError
}

// AddDatatype returns nil if successful and DatatypeDuplicateError when DatatypeDefinition.Identifier is already used.
// returns DatatypeInvalidParentError when using DatatypeDefinition.InheritsFrom if parent datatype does not exist or
// the parent datatype has a different value for DatatypeDefinition.AllowsSubfields
func (d *DatatypeRepository) AddDatatype(definition *DatatypeDefinition) error {
	if _, ok := d.datatypes[definition.Identifier]; ok {
		return DatatypeDuplicateError
	}

	if definition.InheritsFrom != "" {
		parentDatatype, ok := d.datatypes[definition.InheritsFrom]
		if !ok {
			return DatatypeInvalidParentError
		}

		if parentDatatype.definition.AllowsSubfields != definition.AllowsSubfields {
			return DatatypeInvalidParentError
		}
	}

	d.datatypes[definition.Identifier] = &Datatype{definition: *definition}
	return nil
}

// DeleteDatatype deletes a stored datatype inside the repository, returns DatatypeNotFoundError when no datatype is
// found, returns DatatypeIsParentError when trying to delete a datatype that is used by a different datatype in
// DatatypeDefinition.InheritsFrom.
// This function should only be used before starting an application or using the repository to validate data.
func (d *DatatypeRepository) DeleteDatatype(identifier string) error {
	if _, ok := d.datatypes[identifier]; ok {
		for k, v := range d.datatypes {
			if k == identifier {
				continue
			}

			if v.definition.InheritsFrom == identifier {
				return DatatypeIsParentError
			}
		}

		delete(d.datatypes, identifier)
		return nil
	}

	return DatatypeNotFoundError
}
