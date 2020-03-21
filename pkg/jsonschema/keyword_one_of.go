package jsonschema

import (
	"strconv"
)

type OneOf struct{}

var _ Keyword = OneOf{}

func (_ OneOf) Keyword() string {
	return "oneOf"
}

func (_ OneOf) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	var numValid int
	for i, subschema := range input.Schema.JSONValue.([]JSON) {
		childInput := Node{
			Valid:                   true,
			Parent:                  &input,
			Instance:                input.Instance,
			InstanceLocation:        input.InstanceLocation,
			Schema:                  subschema,
			KeywordLocation:         input.KeywordLocation.AddReferenceToken(strconv.Itoa(i)),
			AbsoluteKeywordLocation: input.AbsoluteKeywordLocation.AddReferenceToken(strconv.Itoa(i)),
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

	if numValid != 1 {
		input.Valid = false
	}

	return &input, nil
}
