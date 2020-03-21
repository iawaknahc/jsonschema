package jsonschema

import (
	"strconv"
)

type AdditionalItems struct{}

var _ Keyword = AdditionalItems{}
var _ AnnotatingKeyword = AdditionalItems{}

func (_ AdditionalItems) Keyword() string {
	return "additionalItems"
}

func (_ AdditionalItems) CombineAnnotations(values []interface{}) (interface{}, bool) {
	if len(values) <= 0 {
		return nil, false
	}
	return true, true
}

func (_ AdditionalItems) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}
	a, ok := input.GetAnnotationsFromAdjacentKeywords(Items{})
	if !ok {
		return &input, nil
	}
	if _, ok := a.(bool); ok {
		return &input, nil
	}
	switch j := a.(type) {
	case int:
		applied := false
		for i, item := range arr {
			if i <= j {
				continue
			}
			applied = true
			childInput := Node{
				Valid:                   true,
				Parent:                  &input,
				Instance:                item,
				InstanceLocation:        input.InstanceLocation.AddReferenceToken(strconv.Itoa(i)),
				Schema:                  input.Schema,
				KeywordLocation:         input.KeywordLocation,
				AbsoluteKeywordLocation: input.AbsoluteKeywordLocation,
			}
			child, err := ctx.Apply(childInput)
			if err != nil {
				return nil, err
			}
			if !child.Valid {
				input.Valid = false
			}
			input.Children = append(input.Children, *child)
		}
		if applied {
			input.Annotation = true

		}
	}

	return &input, nil
}
