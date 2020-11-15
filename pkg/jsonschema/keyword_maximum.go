package jsonschema

import (
	"encoding/json"
)

type Maximum struct {
	Maximum json.Number `json:"maximum"`
	Actual  json.Number `json:"actual"`
}

var _ Keyword = Maximum{}

func (_ Maximum) Keyword() string {
	return "maximum"
}

func (_ Maximum) Apply(ctx ApplicationContext, input Node) (*Node, error) {
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
	if numf > limitf {
		input.Valid = false
		input.Info = Maximum{
			Maximum: limit,
			Actual:  num,
		}
	}
	return &input, nil
}
