package jsonschema

import (
	"encoding/json"
)

type Minimum struct {
	Minimum json.Number `json:"minimum"`
	Actual  json.Number `json:"actual"`
}

var _ Keyword = Minimum{}

func (_ Minimum) Keyword() string {
	return "minimum"
}

func (_ Minimum) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	num, ok := input.Instance.(json.Number)
	if !ok {
		return &input, nil
	}
	limit := input.Scope.Schema.JSONValue.(json.Number)
	numf, err := num.Float64()
	if err != nil {
		return nil, err
	}
	limitf, err := limit.Float64()
	if err != nil {
		return nil, err
	}
	if numf < limitf {
		input.Valid = false
		input.Info = Minimum{
			Minimum: limit,
			Actual:  num,
		}
	}
	return &input, nil
}
