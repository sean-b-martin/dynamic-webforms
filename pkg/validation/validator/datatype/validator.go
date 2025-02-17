package datatype

import (
	"encoding/json"
	"github.com/sean-b-martin/dynamic-webforms/pkg/model"
	"github.com/sean-b-martin/dynamic-webforms/pkg/validation/validator/common"
)

type Validator interface {
	New() Validator
	Initialize(id int, rawConstraints *json.RawMessage) common.ValidatorError
	Validate(data *model.WebFormDataRaw) common.ValidatorError
}
