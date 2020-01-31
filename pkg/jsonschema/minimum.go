package jsonschema

import (
	"encoding/json"
)

type Minimum struct {
	Minimum json.Number `json:"minimum"`
	Actual  json.Number `json:"actual"`
}

var _ Keyworder = Minimum{}
var _ Applicator = Minimum{}

func (_ Minimum) Keyword() string {
	return "minimum"
}

func (_ Minimum) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	num, ok := ctx.Instance.(json.Number)
	if !ok {
		return
	}
	limit := ctx.Schema.JSONValue.(json.Number)
	numf, err := num.Float64()
	if err != nil {
		panic(err)
	}
	limitf, err := limit.Float64()
	if err != nil {
		panic(err)
	}
	if numf < limitf {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: Minimum{
				Minimum: limit,
				Actual:  num,
			},
		})
	}
	return
}
