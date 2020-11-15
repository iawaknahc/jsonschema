package jsonschema

type AdditionalProperties struct{}

var _ Keyword = AdditionalProperties{}
var _ AnnotatingKeyword = AdditionalProperties{}

func (_ AdditionalProperties) Keyword() string {
	return "additionalProperties"
}

func (_ AdditionalProperties) CombineAnnotations(values []interface{}) (interface{}, bool) {
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

func (_ AdditionalProperties) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	processedNames := map[string]struct{}{}
	if a, ok := input.GetAnnotationsFromAdjacentKeywords(Properties{}); ok {
		for name := range a.(map[string]struct{}) {
			processedNames[name] = struct{}{}
		}
	}
	if a, ok := input.GetAnnotationsFromAdjacentKeywords(PatternProperties{}); ok {
		for name := range a.(map[string]struct{}) {
			processedNames[name] = struct{}{}
		}
	}

	additionalPropertiesName := map[string]struct{}{}
	for name, val := range obj {
		_, ok := processedNames[name]
		if ok {
			continue
		}
		additionalPropertiesName[name] = struct{}{}
		childInput := Node{
			Valid:            true,
			Parent:           &input,
			Instance:         val,
			InstanceLocation: input.InstanceLocation.AddReferenceToken(name),
			Scope:            input.Scope,
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

	input.Annotation = additionalPropertiesName

	return &input, nil
}
