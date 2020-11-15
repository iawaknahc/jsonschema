package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MaxItems struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MaxItems{}

func (_ MaxItems) Keyword() string {
	return "maxItems"
}

func (_ MaxItems) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}
	maxItems := input.Scope.Schema.JSONValue.(json.Number)
	arrLen := len(arr)
	i, err := strconv.Atoi(string(maxItems))
	if err != nil {
		return nil, err
	}
	if arrLen > i {
		input.Valid = false
		input.Info = MaxItems{
			Expected: maxItems,
			Actual:   arrLen,
		}
	}
	return &input, nil
}
