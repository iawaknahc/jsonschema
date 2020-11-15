package jsonschema

type DependentSchemas struct{}

var _ Keyword = DependentSchemas{}

func (_ DependentSchemas) Keyword() string {
	return "dependentSchemas"
}

func (_ DependentSchemas) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	for name, schema := range input.Scope.Schema.JSONValue.(map[string]JSON) {
		_, ok := obj[name]
		if !ok {
			continue
		}

		childInput := Node{
			Valid:            true,
			Parent:           &input,
			Instance:         input.Instance,
			InstanceLocation: input.InstanceLocation,
			Scope:            input.Scope.Descend(name, schema),
		}

		child, err := ctx.Apply(childInput)
		if err != nil {
			return nil, err
		}

		if !child.Valid {
			input.Valid = false
		}

		input.Children = append(input.Children, *child)
	}

	return &input, nil
}
