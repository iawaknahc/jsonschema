package jsonschema

import (
	"strconv"
)

type AllOf struct{}

var _ Keyword = AllOf{}

func (_ AllOf) Keyword() string {
	return "allOf"
}

func (_ AllOf) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	var numInvalid int
	for i, subschema := range input.Scope.Schema.JSONValue.([]JSON) {
		childInput := Node{
			Valid:            true,
			Parent:           &input,
			Instance:         input.Instance,
			InstanceLocation: input.InstanceLocation,
			Scope:            input.Scope.Descend(strconv.Itoa(i), subschema),
		}
		child, err := ctx.Apply(childInput)
		if err != nil {
			return nil, err
		}

		if !child.Valid {
			numInvalid++
		}

		input.Children = append(input.Children, *child)
	}

	if numInvalid > 0 {
		input.Valid = false
	}

	return &input, nil
}
