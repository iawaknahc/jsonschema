package jsonschema

type Then struct{}

var _ Keyword = Then{}

func (_ Then) Keyword() string {
	return "then"
}

func (_ Then) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	a, ok := input.GetAnnotationsFromAdjacentKeywords(If{})
	if !ok {
		return &input, nil
	}

	if !a.(bool) {
		return &input, nil
	}

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
	input.Valid = child.Valid
	input.Children = append(input.Children, *child)
	return &input, nil
}
