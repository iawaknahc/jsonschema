package jsonschema

import (
	"encoding/json"
	"strconv"
)

type MinItems struct {
	Expected json.Number `json:"expected"`
	Actual   int         `json:"actual"`
}

var _ Keyworder = MinItems{}
var _ Applicator = MinItems{}

func (_ MinItems) Keyword() string {
	return "minItems"
}

func (_ MinItems) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	arr, ok := ctx.Instance.([]interface{})
	if !ok {
		return
	}
	minItems := ctx.Schema.JSONValue.(json.Number)
	arrLen := len(arr)
	i, err := strconv.Atoi(string(minItems))
	if err != nil {
		panic(err)
	}
	if arrLen < i {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: MinItems{
				Expected: minItems,
				Actual:   arrLen,
			},
		})
	}
	return
}
