package jsonschema

type Boolean struct{}

func (b Boolean) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	if valid, ok := ctx.Schema.JSONValue.(bool); ok && !valid {
		errors = append(errors, Error{
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value:                   b,
		})
	}
	return
}
