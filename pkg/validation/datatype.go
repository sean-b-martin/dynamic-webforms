package validation

// DatatypeDefinition contains metadata about a datatype such as the identifier for using the datatype
// and whether the datatype allows subfields or not.
type DatatypeDefinition struct {
	Identifier      string `json:"identifier"`
	DisplayName     string `json:"displayName"`
	AllowsSubfields bool   `json:"allowsSubfields"`
}

// Datatype contains its definition and all rules for validation TODO(add validation rules)
type Datatype struct {
	definition DatatypeDefinition
}

// DatatypeRepository contains the available datatypes.
type DatatypeRepository struct {
	datatypes map[string]*Datatype
}

func NewDatatype(definition DatatypeDefinition) Datatype {
	return Datatype{definition: definition}
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
func (d *DatatypeRepository) AddDatatype(definition *DatatypeDefinition) error {
	if _, ok := d.datatypes[definition.Identifier]; ok {
		return DatatypeDuplicateError
	}

	d.datatypes[definition.Identifier] = &Datatype{definition: *definition}
	return nil
}

// DeleteDatatype deletes a stored datatype inside the repository, returns DatatypeNotFoundError when no datatype is
// found. This function should only be used before starting an application or using the repository to validate data.
func (d *DatatypeRepository) DeleteDatatype(identifier string) error {
	if _, ok := d.datatypes[identifier]; ok {
		delete(d.datatypes, identifier)
		return nil
	}

	return DatatypeNotFoundError
}
