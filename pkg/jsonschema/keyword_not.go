package jsonschema

type Not struct{}

var _ Keyword = Not{}

func (_ Not) Keyword() string {
	return "not"
}

func (_ Not) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	childInput := Node{
		Valid:                   true,
		Parent:                  &input,
		Instance:                input.Instance,
		InstanceLocation:        input.InstanceLocation,
		Schema:                  input.Schema,
		KeywordLocation:         input.KeywordLocation,
		AbsoluteKeywordLocation: input.AbsoluteKeywordLocation,
	}
	child, err := ctx.Apply(childInput)
	if err != nil {
		return nil, err
	}
	input.Valid = !child.Valid
	input.Children = append(input.Children, *child)
	return &input, nil
}
