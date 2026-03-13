package model

import "encoding/json"

type Data struct {
	ID   int             `json:"id"`
	Data json.RawMessage `json:"data"`
}
