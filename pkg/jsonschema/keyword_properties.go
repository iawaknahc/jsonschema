package jsonschema

type Properties struct{}

var _ Keyword = Properties{}
var _ AnnotatingKeyword = Properties{}

func (_ Properties) Keyword() string {
	return "properties"
}

func (_ Properties) CombineAnnotations(values []interface{}) (interface{}, bool) {
	if len(values) <= 0 {
		return nil, false
	}

	merged := map[string]struct{}{}
	for _, v := range values {
		for name := range v.(map[string]struct{}) {
			merged[name] = struct{}{}
		}
	}

	return merged, true
}

func (_ Properties) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	propertiesName := map[string]struct{}{}
	for name, schema := range input.Scope.Schema.JSONValue.(map[string]JSON) {
		if val, ok := obj[name]; ok {
			propertiesName[name] = struct{}{}
			childInput := Node{
				Valid:            true,
				Parent:           &input,
				Instance:         val,
				InstanceLocation: input.InstanceLocation.AddReferenceToken(name),
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
	}

	input.Annotation = propertiesName

	return &input, nil
}
