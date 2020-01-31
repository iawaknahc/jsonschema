package jsonschema

import (
	"strconv"
)

type AnyOf struct {
	Result []bool `json:"result"`
}

var _ Keyworder = AnyOf{}
var _ Applicator = AnyOf{}

func (_ AnyOf) Keyword() string {
	return "anyOf"
}

func (_ AnyOf) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	var result []bool
	var numValid int
	var numInvalid int
	for i, subschema := range ctx.Schema.JSONValue.([]JSON) {
		c := ctx
		c.Schema = subschema
		c.KeywordLocation = c.KeywordLocation.AddReferenceToken(strconv.Itoa(i))
		c.AbsoluteKeywordLocation = c.AbsoluteKeywordLocation.AddReferenceToken(strconv.Itoa(i))
		a, e := c.Apply()
		if len(e) > 0 {
			numInvalid++
			result = append(result, false)
		} else {
			numValid++
			result = append(result, true)
			annotations = append(annotations, a...)
		}
		errors = append(errors, e...)
	}

	if numValid <= 0 {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: AnyOf{
				Result: result,
			},
		})
	} else {
		errors = nil
	}

	if len(errors) > 0 {
		annotations = nil
	}

	return
}
