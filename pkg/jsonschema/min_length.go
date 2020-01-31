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

var _ Keyworder = MinLength{}
var _ Applicator = MinLength{}

func (_ MinLength) Keyword() string {
	return "minLength"
}

func (_ MinLength) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	str, ok := ctx.Instance.(string)
	if !ok {
		return
	}
	minItems := ctx.Schema.JSONValue.(json.Number)
	length := utf8.RuneCountInString(str)
	i, err := strconv.Atoi(string(minItems))
	if err != nil {
		panic(err)
	}
	if length < i {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: MinLength{
				Expected: minItems,
				Actual:   length,
			},
		})
	}
	return
}
