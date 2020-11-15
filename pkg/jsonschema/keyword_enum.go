package jsonschema

import (
	"reflect"
)

type Enum struct {
	Expected []interface{} `json:"expected"`
	Actual   interface{}   `json:"actual"`
}

var _ Keyword = Enum{}

func (_ Enum) Keyword() string {
	return "enum"
}

func (_ Enum) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	enumValues := ToFloat64(UnwrapJSON(input.Scope.Schema)).([]interface{})
	value := ToFloat64(input.Instance)

	eq := false
	for _, enumValue := range enumValues {
		if reflect.DeepEqual(value, enumValue) {
			eq = true
			break
		}
	}

	if !eq {
		input.Valid = false
		input.Info = Enum{
			Expected: enumValues,
			Actual:   value,
		}
	}

	return &input, nil
}
