package jsonschema

type Else struct{}

var _ Keyword = Else{}

func (_ Else) Keyword() string {
	return "else"
}

func (_ Else) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	a, ok := input.GetAnnotationsFromAdjacentKeywords(If{})
	if !ok {
		return &input, nil
	}

	if a.(bool) {
		return &input, nil
	}

	childInput := Node{
		Valid:            true,
		Parent:           &input,
		Instance:         input.Instance,
		InstanceLocation: input.InstanceLocation,
		Scope:            input.Scope,
	}
	child, err := ctx.Apply(childInput)
	if err != nil {
		return nil, err
	}
	input.Valid = child.Valid
	input.Children = append(input.Children, *child)
	return &input, nil
}
