package jsonschema

import (
	"encoding/json"
	"strconv"
	"unicode/utf8"
)

type MaxLength struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyword = MaxLength{}

func (_ MaxLength) Keyword() string {
	return "maxLength"
}

func (_ MaxLength) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	str, ok := input.Instance.(string)
	if !ok {
		return &input, nil
	}
	maxItems := input.Scope.Schema.JSONValue.(json.Number)
	length := utf8.RuneCountInString(str)
	i, err := strconv.Atoi(string(maxItems))
	if err != nil {
		return nil, err
	}
	if length > i {
		input.Valid = false
		input.Info = MaxLength{
			Expected: maxItems,
			Actual:   length,
		}
	}
	return &input, nil
}
