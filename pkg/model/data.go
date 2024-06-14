package model

import "encoding/json"

// WebFormDataRaw contains the ID of the schema element that is used to validate Data. Data gets converted to concrete
// datatypes when necessary. WebFormDataRaw is used for validation in order to prevent errors like integer overflows.
type WebFormDataRaw struct {
	SchemaElementID int               `json:"schemaElementID"`
	Data            []json.RawMessage `json:"data"`
}

// WebFormData can be used when Data from WebFormDataRaw was successfully validated. Original Data is not validated here
// in order to prevent errors like integer overflows when deserializing json.
type WebFormData struct {
	SchemaElementID int          `json:"schemaElementID"`
	Data            []ParsedData `json:"data"`
}

type ParsedData interface{}
