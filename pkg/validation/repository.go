package validation

import (
	"fmt"
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator"
)

type BaseValidatorRepository interface {
	AddValidator(name string, validator validator.Factory) error
}

type ValidatorRepository interface {
	BaseValidatorRepository
	GetValidator(name string) (validator.Factory, error)
}

type GlobalValidatorRepository interface {
	BaseValidatorRepository
	AddDatatype(datatype ValidatorDatatype) error
	AddToDatatype(datatype string, name string, validator validator.Factory) error
	GetValidator(datatype string, name string) (validator.Factory, error)
	GetDatatype(name string) (Datatype, error)
}

type validatorRepository struct {
	validators map[string]validator.Factory
}

func (v *validatorRepository) AddValidator(name string, validator validator.Factory) error {
	if _, ok := v.validators[name]; ok {
		return fmt.Errorf("validator already exists: %s", name)
	}
	v.validators[name] = validator
	return nil
}

func (v *validatorRepository) GetValidator(name string) (validator.Factory, error) {
	val, ok := v.validators[name]
	if !ok {
		return nil, fmt.Errorf("validator not found: %s", name)
	}
	return val, nil
}

type globalValidatorRepository struct {
	repository ValidatorRepository
	datatypes  map[string]ValidatorDatatype
}

func (g *globalValidatorRepository) AddValidator(name string, validator validator.Factory) error {
	return g.repository.AddValidator(name, validator)
}

func (g *globalValidatorRepository) GetValidator(datatype string, name string) (validator.Factory, error) {
	// global validator
	val, err := g.repository.GetValidator(name)
	if err == nil {
		return val, nil
	}

	// datatype validator
	validatorType, ok := g.datatypes[datatype]
	if !ok {
		return nil, fmt.Errorf("datatype not found: %s", datatype)
	}
	return validatorType.GetValidator(name)
}

func (g *globalValidatorRepository) AddDatatype(datatype ValidatorDatatype) error {
	name := datatype.GetDefinition().Name
	if _, ok := g.datatypes[name]; ok {
		return fmt.Errorf("datatype already exists: %s", name)
	}
	g.datatypes[name] = datatype
	return nil
}

func (g *globalValidatorRepository) AddToDatatype(datatype string, name string, validator validator.Factory) error {
	found, ok := g.datatypes[datatype]
	if !ok {
		return fmt.Errorf("datatype not found: %s", datatype)
	}
	return found.AddValidator(name, validator)
}

func (g *globalValidatorRepository) GetDatatype(name string) (Datatype, error) {
	val, ok := g.datatypes[name]
	if !ok {
		return nil, fmt.Errorf("validator not found: %s", name)
	}
	return val, nil
}

type ValidatorRepoOption func(*validatorRepository)

func NewValidatorRepository(opts ...ValidatorRepoOption) ValidatorRepository {
	return ValidatorRepository(&validatorRepository{validators: make(map[string]validator.Factory)})
}

func WithValidators(validators map[string]validator.Factory) ValidatorRepoOption {
	return func(r *validatorRepository) {
		for name, factory := range validators {
			r.validators[name] = factory
		}
	}
}

func WithValidator(name string, factory validator.Factory) ValidatorRepoOption {
	return func(r *validatorRepository) {
		r.validators[name] = factory
	}
}

func NewGlobalValidatorRepository() GlobalValidatorRepository {
	return GlobalValidatorRepository(&globalValidatorRepository{repository: NewValidatorRepository(), datatypes: make(map[string]ValidatorDatatype)})
}
