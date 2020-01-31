package jsonschema

import (
	"strconv"
)

type AdditionalItems struct{}

var _ Keyworder = AdditionalItems{}
var _ Applicator = AdditionalItems{}
var _ Annotator = AdditionalItems{}

func (_ AdditionalItems) Keyword() string {
	return "additionalItems"
}

func (_ AdditionalItems) MergeAnnotations(annotations []Annotation) (*Annotation, bool) {
	if len(annotations) <= 0 {
		return nil, false
	}
	out := annotations[0]
	return &out, true
}

func (_ AdditionalItems) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	arr, ok := ctx.Instance.([]interface{})
	if !ok {
		return
	}
	a, ok := ctx.GetAnnotation(Items{})
	if !ok {
		return
	}
	if _, ok := a.Value.(bool); ok {
		return
	}
	switch j := a.Value.(type) {
	case int:
		for i, item := range arr {
			if i <= j {
				continue
			}
			c := ctx
			c.Instance = item
			c.InstanceLocation = c.InstanceLocation.AddReferenceToken(strconv.Itoa(i))
			childA, childE := c.Apply()
			annotations = append(annotations, childA...)
			errors = append(errors, childE...)
		}
		annotations = append(annotations, Annotation{
			InstanceLocation:        ctx.InstanceLocation,
			Keyword:                 ctx.Keyword,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value:                   true,
		})
	}

	return
}
