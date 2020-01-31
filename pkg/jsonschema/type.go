package jsonschema

import (
	"encoding/json"
	"strings"
)

type Type struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
}

var _ Keyworder = Type{}
var _ Applicator = Type{}

func (_ Type) Keyword() string {
	return "type"
}

func (_ Type) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	// Prepare expected.
	var expected []string
	switch a := UnwrapJSON(ctx.Schema).(type) {
	case string:
		expected = []string{a}
	case []interface{}:
		for _, b := range a {
			if c, ok := b.(string); ok {
				expected = append(expected, c)
			}
		}
	}

	// Prepare actual
	var actual []string
	switch d := ctx.Instance.(type) {
	case nil:
		actual = append(actual, "null")
	case string:
		actual = append(actual, "string")
	case bool:
		actual = append(actual, "boolean")
	case map[string]interface{}:
		actual = append(actual, "object")
	case []interface{}:
		actual = append(actual, "array")
	case json.Number:
		if strings.ContainsRune(string(d), '.') {
			actual = append(actual, "number")
		} else {
			actual = append(actual, "number", "integer")
		}
	}

	intersection := IntersectString(expected, actual)
	if len(intersection) <= 0 {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: Type{
				Expected: expected,
				Actual:   actual,
			},
		})
	}
	return
}

func IntersectString(a []string, b []string) (out []string) {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				out = append(out, x)
			}
		}
	}
	return
}
