package jsonschema

import (
	"encoding/json"
)

type ExclusiveMinimum struct {
	ExclusiveMinimum json.Number `json:"exclusiveMinimum"`
	Actual           json.Number `json:"actual"`
}

var _ Keyworder = ExclusiveMinimum{}
var _ Applicator = ExclusiveMinimum{}

func (_ ExclusiveMinimum) Keyword() string {
	return "exclusiveMinimum"
}

func (_ ExclusiveMinimum) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
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
	if numf <= limitf {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: ExclusiveMinimum{
				ExclusiveMinimum: limit,
				Actual:           num,
			},
		})
	}
	return
}
