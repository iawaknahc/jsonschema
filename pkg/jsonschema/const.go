package jsonschema

import (
	"reflect"
)

type Const struct {
	Const interface{} `json:"const"`
}

var _ Keyworder = Const{}
var _ Applicator = Const{}

func (_ Const) Keyword() string {
	return "const"
}

func (_ Const) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	constValue := ToFloat64(UnwrapJSON(ctx.Schema))
	value := ToFloat64(ctx.Instance)

	if !reflect.DeepEqual(value, constValue) {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: Const{
				Const: constValue,
			},
		})
	}
	return
}
