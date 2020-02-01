package jsonschema

import (
	"strconv"
)

type Items struct{}

var _ Keyword = Items{}
var _ AnnotatingKeyword = Items{}

func (_ Items) Keyword() string {
	return "items"
}

func (_ Items) CombineAnnotations(values []interface{}) (interface{}, bool) {
	if len(values) <= 0 {
		return nil, false
	}

	largestIndex := -1
	for _, a := range values {
		switch v := a.(type) {
		case bool:
			return true, true
		case int:
			if v > largestIndex {
				largestIndex = v
			}
		}
	}

	return largestIndex, true
}

func (_ Items) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}

	largestIndex := -1
	switch schema := input.Schema.JSONValue.(type) {
	case []JSON:
		for i := 0; i < len(arr) && i < len(schema); i++ {
			largestIndex = i
			item := arr[i]
			childInput := Node{
				Valid:                   true,
				Parent:                  &input,
				Instance:                item,
				InstanceLocation:        input.InstanceLocation.AddReferenceToken(strconv.Itoa(i)),
				Schema:                  schema[i],
				KeywordLocation:         input.KeywordLocation.AddReferenceToken(strconv.Itoa(i)),
				AbsoluteKeywordLocation: input.AbsoluteKeywordLocation.AddReferenceToken(strconv.Itoa(i)),
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
	default:
		for i, item := range arr {
			largestIndex = i
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
	}

	if largestIndex == len(arr)-1 {
		input.Annotation = true
	} else {
		input.Annotation = largestIndex
	}

	return &input, nil
}
