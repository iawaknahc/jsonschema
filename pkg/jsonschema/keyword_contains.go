package jsonschema

import (
	"strconv"
)

type Contains struct{}

var _ Keyword = Contains{}

func (_ Contains) Keyword() string {
	return "contains"
}

func (_ Contains) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}
	numValid := 0
	for i := 0; i < len(arr); i++ {
		item := arr[i]
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
		if child.Valid {
			numValid++
		}
		input.Children = append(input.Children, *child)
	}

	if numValid <= 0 {
		input.Valid = false
	}

	return &input, nil
}
