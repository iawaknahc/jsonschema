package jsonschema

type Format struct{}

var _ Keyword = Format{}

func (_ Format) Keyword() string {
	return "format"
}

func (_ Format) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	format := input.Scope.Schema.JSONValue.(string)
	input.Annotation = format

	checker, ok := ctx.Collection.FormatChecker[format]
	if !ok {
		return &input, nil
	}

	err := checker.CheckFormat(input.Instance)
	if err != nil {
		input.Valid = false
		input.Info = err
	}

	return &input, nil
}
