package jsonschema

import (
	"encoding/json"
)

type MultipleOf struct {
	Dividend float64 `json:"dividend"`
	Divisor  float64 `json:"divisor"`
	Quotient float64 `json:"quotient"`
}

var _ Keyword = MultipleOf{}

func (_ MultipleOf) Keyword() string {
	return "multipleOf"
}

func (_ MultipleOf) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	dividendStr, ok := input.Instance.(json.Number)
	if !ok {
		return &input, nil
	}
	dividend, err := dividendStr.Float64()
	if err != nil {
		return nil, err
	}
	divisor, err := input.Schema.JSONValue.(json.Number).Float64()
	if err != nil {
		return nil, err
	}
	quotient := dividend / divisor
	ok = float64(int(quotient)) == quotient
	if !ok {
		input.Valid = false
		input.Info = MultipleOf{
			Dividend: dividend,
			Divisor:  divisor,
			Quotient: quotient,
		}
	}
	return &input, nil
}
