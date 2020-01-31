package jsonschema

type Not struct{}

var _ Keyworder = Not{}
var _ Applicator = Not{}

func (_ Not) Keyword() string {
	return "not"
}

func (_ Not) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	_, e := ctx.Apply()
	if len(e) <= 0 {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value:                   Not{},
		})
	}
	return
}
