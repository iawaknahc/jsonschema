package jsonschema

import (
	"strconv"
)

type AnyOf struct{}

var _ Keyword = AnyOf{}

func (_ AnyOf) Keyword() string {
	return "anyOf"
}

func (_ AnyOf) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	var numValid int
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
