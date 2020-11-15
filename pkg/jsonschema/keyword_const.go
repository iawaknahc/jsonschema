package jsonschema

import (
	"reflect"
)

type Const struct {
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}

var _ Keyword = Const{}

func (_ Const) Keyword() string {
	return "const"
}

func (_ Const) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	constValue := ToFloat64(UnwrapJSON(input.Scope.Schema))
	value := ToFloat64(input.Instance)

	if !reflect.DeepEqual(value, constValue) {
		input.Valid = false
		input.Info = Const{
			Expected: constValue,
			Actual:   value,
		}
	}

	return &input, nil
}
