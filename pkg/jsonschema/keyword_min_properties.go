package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MinProperties struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MinProperties{}

func (_ MinProperties) Keyword() string {
	return "minProperties"
}

func (_ MinProperties) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}
	limit := input.Scope.Schema.JSONValue.(json.Number)
	arrLen := len(obj)
	i, err := strconv.Atoi(string(limit))
	if err != nil {
		return nil, err
	}
	if arrLen < i {
		input.Valid = false
		input.Info = MinProperties{
			Expected: limit,
			Actual:   arrLen,
		}
	}
	return &input, nil
}
