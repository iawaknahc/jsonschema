package jsonschema

import (
	"encoding/json"
)

type ExclusiveMaximum struct {
	ExclusiveMaximum json.Number `json:"exclusiveMaximum"`
	Actual           json.Number `json:"actual"`
}

var _ Keyworder = ExclusiveMaximum{}
var _ Applicator = ExclusiveMaximum{}

func (_ ExclusiveMaximum) Keyword() string {
	return "exclusiveMaximum"
}

func (_ ExclusiveMaximum) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
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
	if numf >= limitf {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: ExclusiveMaximum{
				ExclusiveMaximum: limit,
				Actual:           num,
			},
		})
	}
	return
}
