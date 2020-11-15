package jsonschema

import (
	"encoding/json"
	"strconv"
)

type Contains struct {
	Min    *int `json:"min,omitempty"`
	Max    *int `json:"max,omitempty"`
	Actual int  `json:"actual"`
}

var _ Keyword = Contains{}

func (_ Contains) Keyword() string {
	return "contains"
}

func (_ Contains) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}

	var min *int
	var max *int

	parentSchema := input.Parent.Scope.Schema.JSONValue.(map[string]JSON)

	if minContains, ok := parentSchema["minContains"].JSONValue.(json.Number); ok {
		i, err := strconv.Atoi(string(minContains))
		if err != nil {
			return nil, err
		}
		min = &i
	}

	if maxContains, ok := parentSchema["maxContains"].JSONValue.(json.Number); ok {
		i, err := strconv.Atoi(string(maxContains))
		if err != nil {
			return nil, err
		}
		max = &i
	}

	actual := 0
	for i := 0; i < len(arr); i++ {
		item := arr[i]
		childInput := Node{
			Valid:            true,
			Parent:           &input,
			Instance:         item,
			InstanceLocation: input.InstanceLocation.AddReferenceToken(strconv.Itoa(i)),
			Scope:            input.Scope,
		}
		child, err := ctx.Apply(childInput)
		if err != nil {
			return nil, err
		}
		if child.Valid {
			actual++
		}
		input.Children = append(input.Children, *child)
	}

	info := Contains{
		Min:    min,
		Max:    max,
		Actual: actual,
	}

	if min == nil {
		if actual < 1 {
			input.Valid = false
			input.Info = info
		}
	} else {
		if actual < *min {
			input.Valid = false
			input.Info = info
		}
	}
	if max != nil {
		if actual > *max {
			input.Valid = false
			input.Info = info
		}
	}

	return &input, nil
}
