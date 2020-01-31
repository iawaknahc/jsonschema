package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MaxItems struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyworder = MaxItems{}
var _ Applicator = MaxItems{}

func (_ MaxItems) Keyword() string {
	return "maxItems"
}

func (_ MaxItems) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	arr, ok := ctx.Instance.([]interface{})
	if !ok {
		return
	}
	maxItems := ctx.Schema.JSONValue.(json.Number)
	arrLen := len(arr)
	i, err := strconv.Atoi(string(maxItems))
	if err != nil {
		panic(err)
	}
	if arrLen > i {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: MaxItems{
				Expected: maxItems,
				Actual:   arrLen,
			},
		})
	}
	return
}
