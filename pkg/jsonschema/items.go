package jsonschema

import (
	"strconv"
)

type Items struct{}

var _ Keyworder = Items{}
var _ Applicator = Items{}
var _ Annotator = Items{}

func (_ Items) Keyword() string {
	return "items"
}

func (_ Items) MergeAnnotations(annotations []Annotation) (*Annotation, bool) {
	if len(annotations) <= 0 {
		return nil, false
	}

	out := annotations[0]
	hasTrue := false
	largestIndex := -1

	for _, a := range annotations {
		switch v := a.Value.(type) {
		case bool:
			hasTrue = true
		case int:
			if v > largestIndex {
				largestIndex = v
			}
		}
	}

	if hasTrue {
		out.Value = true
	} else {
		out.Value = largestIndex
	}

	return &out, true
}

func (_ Items) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	arr, ok := ctx.Instance.([]interface{})
	if !ok {
		return
	}

	largestIndex := -1
	switch schema := ctx.Schema.JSONValue.(type) {
	case []JSON:
		for i := 0; i < len(arr) && i < len(schema); i++ {
			largestIndex = i
			item := arr[i]
			c := ctx
			c.Schema = schema[i]
			c.KeywordLocation = c.KeywordLocation.AddReferenceToken(strconv.Itoa(i))
			c.AbsoluteKeywordLocation = c.AbsoluteKeywordLocation.AddReferenceToken(strconv.Itoa(i))
			c.Instance = item
			c.InstanceLocation = c.InstanceLocation.AddReferenceToken(strconv.Itoa(i))
			childA, childE := c.Apply()
			annotations = append(annotations, childA...)
			errors = append(errors, childE...)
		}
	default:
		for i, item := range arr {
			largestIndex = i
			c := ctx
			c.Instance = item
			c.InstanceLocation = c.InstanceLocation.AddReferenceToken(strconv.Itoa(i))
			childA, childE := c.Apply()
			annotations = append(annotations, childA...)
			errors = append(errors, childE...)
		}
	}

	var value interface{}
	if largestIndex == len(arr)-1 {
		value = true
	} else {
		value = largestIndex
	}
	annotations = append(annotations, Annotation{
		InstanceLocation:        ctx.InstanceLocation,
		Keyword:                 ctx.Keyword,
		KeywordLocation:         ctx.KeywordLocation,
		AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
		Value:                   value,
	})

	return
}
