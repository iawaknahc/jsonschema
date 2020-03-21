package jsonschema

import (
	"encoding/json"
)

type ExclusiveMinimum struct {
	ExclusiveMinimum json.Number `json:"exclusiveMinimum"`
	Actual           json.Number `json:"actual"`
}

var _ Keyword = ExclusiveMinimum{}

func (_ ExclusiveMinimum) Keyword() string {
	return "exclusiveMinimum"
}

func (_ ExclusiveMinimum) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	num, ok := input.Instance.(json.Number)
	if !ok {
		return &input, nil
	}
	limit := input.Schema.JSONValue.(json.Number)
	numf, err := num.Float64()
	if err != nil {
		return nil, err
	}
	limitf, err := limit.Float64()
	if err != nil {
		return nil, err
	}
	if numf <= limitf {
		input.Valid = false
		input.Info = ExclusiveMinimum{
			ExclusiveMinimum: limit,
			Actual:           num,
		}
	}
	return &input, nil
}
