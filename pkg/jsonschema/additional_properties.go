package jsonschema

type AdditionalProperties struct{}

var _ Keyworder = AdditionalProperties{}
var _ Applicator = AdditionalProperties{}

func (_ AdditionalProperties) Keyword() string {
	return "additionalProperties"
}

func (_ AdditionalProperties) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	obj, ok := ctx.Instance.(map[string]interface{})
	if !ok {
		return
	}

	processedNames := map[string]struct{}{}
	if a, ok := ctx.GetAnnotation(Properties{}); ok {
		for name := range a.Value.(map[string]struct{}) {
			processedNames[name] = struct{}{}
		}
	}
	if a, ok := ctx.GetAnnotation(PatternProperties{}); ok {
		for name := range a.Value.(map[string]struct{}) {
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
		c := ctx
		c.Instance = val
		c.InstanceLocation = c.InstanceLocation.AddReferenceToken(name)
		childA, childE := c.Apply()
		annotations = append(annotations, childA...)
		errors = append(errors, childE...)
	}

	annotations = append(annotations, Annotation{
		InstanceLocation:        ctx.InstanceLocation,
		Keyword:                 ctx.Keyword,
		KeywordLocation:         ctx.KeywordLocation,
		AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
		Value:                   additionalPropertiesName,
	})

	return
}
