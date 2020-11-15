package jsonschema

import (
	"encoding/json"
	"strconv"
	"unicode/utf8"
)

type MinLength struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MinLength{}

func (_ MinLength) Keyword() string {
	return "minLength"
}

func (_ MinLength) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	str, ok := input.Instance.(string)
	if !ok {
		return &input, nil
	}
	minItems := input.Scope.Schema.JSONValue.(json.Number)
	length := utf8.RuneCountInString(str)
	i, err := strconv.Atoi(string(minItems))
	if err != nil {
		return nil, err
	}
	if length < i {
		input.Valid = false
		input.Info = MinLength{
			Expected: minItems,
			Actual:   length,
		}
	}
	return &input, nil
}
