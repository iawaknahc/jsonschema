package jsonschema

type PropertyNames struct{}

var _ Keyword = PropertyNames{}

func (_ PropertyNames) Keyword() string {
	return "propertyNames"
}

func (_ PropertyNames) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}
	for name := range obj {
		childInput := Node{
			Valid:                   true,
			Parent:                  &input,
			Instance:                name,
			InstanceLocation:        input.InstanceLocation.AddReferenceToken(name),
			Schema:                  input.Schema,
			KeywordLocation:         input.KeywordLocation,
			AbsoluteKeywordLocation: input.AbsoluteKeywordLocation,
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
