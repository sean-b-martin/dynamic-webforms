package validator

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
)

type DatatypeValidator interface {
	New() DatatypeValidator
	Initialize(id int, rawConstraints *json.RawMessage) FieldValidatorError
	ValidateFieldData(data *model.WebFormDataRaw) FieldValidatorError
}
