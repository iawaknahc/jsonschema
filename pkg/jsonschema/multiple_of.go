package jsonschema

import (
	"encoding/json"
)

type MultipleOf struct {
	Dividend float64 `json:"dividend"`
	Divisor  float64 `json:"divisor"`
	Quotient float64 `json:"quotient"`
}

var _ Keyworder = MultipleOf{}
var _ Applicator = MultipleOf{}

func (_ MultipleOf) Keyword() string {
	return "multipleOf"
}

func (_ MultipleOf) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	dividendStr, ok := ctx.Instance.(json.Number)
	if !ok {
		return
	}
	dividend, err := dividendStr.Float64()
	if err != nil {
		panic(err)
	}
	divisor, err := ctx.Schema.JSONValue.(json.Number).Float64()
	if err != nil {
		panic(err)
	}
	quotient := dividend / divisor
	ok = float64(int(quotient)) == quotient
	if !ok {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: MultipleOf{
				Dividend: dividend,
				Divisor:  divisor,
				Quotient: quotient,
			},
		})
	}
	return
}
