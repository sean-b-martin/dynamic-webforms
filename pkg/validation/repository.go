package validation

// DatatypeRepository contains the available datatypes.
type DatatypeRepository struct {
	datatypes map[string]*Datatype
}
type DatatypeRepositoryOption func(*DatatypeRepository) error

func NewDatatypeRepository(options ...DatatypeRepositoryOption) (DatatypeRepository, error) {
	repository := DatatypeRepository{datatypes: make(map[string]*Datatype)}

	for _, option := range options {
		if err := option(&repository); err != nil {
			return repository, err
		}
	}

	return repository, nil
}

func WithDefaultDatatypes() DatatypeRepositoryOption {
	return func(repository *DatatypeRepository) error {
		defaultDatatypes := []*Datatype{nil}
		for _, datatype := range defaultDatatypes {
			err := repository.AddDatatype(datatype)
			if err != nil {
				return err
			}
		}
		return nil
	}
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
// returns Datatype,nil if successful and nil, ErrDatatypeNotFound when no datatype was found
func (d *DatatypeRepository) GetDatatype(id string) (*Datatype, error) {
	if datatype, ok := d.datatypes[id]; ok {
		return datatype, nil
	}

	return nil, ErrDatatypeNotFound
}

// AddDatatype returns nil if successful and ErrDatatypeDuplicate when DatatypeDefinition.ID is already used.
func (d *DatatypeRepository) AddDatatype(datatype *Datatype) error {
	if _, ok := d.datatypes[datatype.definition.ID]; ok {
		return ErrDatatypeDuplicate
	}
	d.datatypes[datatype.definition.ID] = datatype
	return nil
}

// DeleteDatatype deletes a stored datatype inside the repository, returns ErrDatatypeNotFound when no datatype is
// found. This function should only be used before starting an application or using the repository to validate data.
func (d *DatatypeRepository) DeleteDatatype(id string) error {
	if _, ok := d.datatypes[id]; ok {
		delete(d.datatypes, id)
		return nil
	}

	return ErrDatatypeNotFound
}
