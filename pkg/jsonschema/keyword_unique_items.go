package jsonschema

import (
	"reflect"
)

type UniqueItems struct{}

var _ Keyword = UniqueItems{}

func (_ UniqueItems) Keyword() string {
	return "uniqueItems"
}

func (_ UniqueItems) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	arr, ok := input.Instance.([]interface{})
	if !ok {
		return &input, nil
	}
	unique := input.Scope.Schema.JSONValue.(bool)
	if !unique {
		return &input, nil
	}

	arr = ToFloat64(arr).([]interface{})
loop:
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			a := arr[i]
			b := arr[j]
			if reflect.DeepEqual(a, b) {
				input.Valid = false
				input.Info = UniqueItems{}
				break loop
			}
		}
	}

	return &input, nil
}
