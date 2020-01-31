package jsonschema

import (
	"strconv"
)

type AllOf struct {
	Result []bool `json:"result"`
}

var _ Keyworder = AllOf{}
var _ Applicator = AllOf{}

func (_ AllOf) Keyword() string {
	return "allOf"
}

func (_ AllOf) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
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

	if numInvalid > 0 {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: AllOf{
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
