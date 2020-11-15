package jsonschema

import (
	"encoding/json"
)

type ExclusiveMaximum struct {
	ExclusiveMaximum json.Number `json:"exclusiveMaximum"`
	Actual           json.Number `json:"actual"`
}

var _ Keyword = ExclusiveMaximum{}

func (_ ExclusiveMaximum) Keyword() string {
	return "exclusiveMaximum"
}

func (_ ExclusiveMaximum) Apply(ctx ApplicationContext, input Node) (*Node, error) {
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
	if numf >= limitf {
		input.Valid = false
		input.Info = ExclusiveMaximum{
			ExclusiveMaximum: limit,
			Actual:           num,
		}
	}
	return &input, nil
}
