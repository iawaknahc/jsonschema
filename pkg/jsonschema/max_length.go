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

var _ Keyworder = MaxLength{}
var _ Applicator = MaxLength{}

func (_ MaxLength) Keyword() string {
	return "maxLength"
}

func (_ MaxLength) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	str, ok := ctx.Instance.(string)
	if !ok {
		return
	}
	maxItems := ctx.Schema.JSONValue.(json.Number)
	length := utf8.RuneCountInString(str)
	i, err := strconv.Atoi(string(maxItems))
	if err != nil {
		panic(err)
	}
	if length > i {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: MaxLength{
				Expected: maxItems,
				Actual:   length,
			},
		})
	}
	return
}
