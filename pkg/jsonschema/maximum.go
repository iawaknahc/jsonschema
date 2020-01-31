package jsonschema

import (
	"encoding/json"
)

type Maximum struct {
	Maximum json.Number `json:"maximum"`
	Actual  json.Number `json:"actual"`
}

var _ Keyworder = Maximum{}
var _ Applicator = Maximum{}

func (_ Maximum) Keyword() string {
	return "maximum"
}

func (_ Maximum) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
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
	if numf > limitf {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: Maximum{
				Maximum: limit,
				Actual:  num,
			},
		})
	}
	return
}
