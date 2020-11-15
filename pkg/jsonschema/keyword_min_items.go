package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MinItems struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MinItems{}

func (_ MinItems) Keyword() string {
	return "minItems"
}

func (_ MinItems) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}
	minItems := input.Scope.Schema.JSONValue.(json.Number)
	arrLen := len(arr)
	i, err := strconv.Atoi(string(minItems))
	if err != nil {
		return nil, err
	}
	if arrLen < i {
		input.Valid = false
		input.Info = MinItems{
			Expected: minItems,
			Actual:   arrLen,
		}
	}
	return &input, nil
}
