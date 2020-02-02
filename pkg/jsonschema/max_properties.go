package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MaxProperties struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MaxProperties{}

func (_ MaxProperties) Keyword() string {
	return "maxProperties"
}

func (_ MaxProperties) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}
	limit := input.Schema.JSONValue.(json.Number)
	arrLen := len(obj)
	i, err := strconv.Atoi(string(limit))
	if err != nil {
		return nil, err
	}
	if arrLen > i {
		input.Valid = false
		input.Info = MaxProperties{
			Expected: limit,
			Actual:   arrLen,
		}
	}
	return &input, nil
}
