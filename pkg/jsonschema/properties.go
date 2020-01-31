package jsonschema

type Properties struct{}

var _ Keyworder = Properties{}
var _ Applicator = Properties{}
var _ Annotator = Properties{}

func (_ Properties) Keyword() string {
	return "properties"
}

func (_ Properties) MergeAnnotations(annotations []Annotation) (*Annotation, bool) {
	if len(annotations) <= 0 {
		return nil, false
	}

	out := annotations[0]
	merged := map[string]struct{}{}
	for _, a := range annotations {
		for name := range a.Value.(map[string]struct{}) {
			merged[name] = struct{}{}
		}
	}

	out.Value = merged
	return &out, true
}

func (_ Properties) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	obj, ok := ctx.Instance.(map[string]interface{})
	if !ok {
		return
	}

	propertiesName := map[string]struct{}{}
	for name, schema := range ctx.Schema.JSONValue.(map[string]JSON) {
		if val, ok := obj[name]; ok {
			propertiesName[name] = struct{}{}
			c := ctx
			c.Schema = schema
			c.KeywordLocation = c.KeywordLocation.AddReferenceToken(name)
			c.AbsoluteKeywordLocation = c.AbsoluteKeywordLocation.AddReferenceToken(name)
			c.Instance = val
			c.InstanceLocation = c.InstanceLocation.AddReferenceToken(name)
			childA, childE := c.Apply()
			annotations = append(annotations, childA...)
			errors = append(errors, childE...)
		}
	}

	// TODO: Add error for properties

	annotations = append(annotations, Annotation{
		InstanceLocation:        ctx.InstanceLocation,
		Keyword:                 ctx.Keyword,
		KeywordLocation:         ctx.KeywordLocation,
		AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
		Value:                   propertiesName,
	})

	return
}
