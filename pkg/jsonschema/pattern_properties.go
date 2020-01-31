package jsonschema

import (
	"regexp"
)

type PatternProperties struct{}

var _ Keyworder = PatternProperties{}
var _ Applicator = PatternProperties{}
var _ Annotator = PatternProperties{}

func (_ PatternProperties) Keyword() string {
	return "patternProperties"
}

func (_ PatternProperties) MergeAnnotations(annotations []Annotation) (*Annotation, bool) {
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

func (_ PatternProperties) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	obj, ok := ctx.Instance.(map[string]interface{})
	if !ok {
		return
	}

	patternPropertiesName := map[string]struct{}{}
	for pattern, schema := range ctx.Schema.JSONValue.(map[string]JSON) {
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		for name, val := range obj {
			if re.MatchString(name) {
				patternPropertiesName[name] = struct{}{}
				c := ctx
				c.Schema = schema
				c.KeywordLocation = c.KeywordLocation.AddReferenceToken(pattern)
				c.AbsoluteKeywordLocation = c.AbsoluteKeywordLocation.AddReferenceToken(pattern)
				c.Instance = val
				c.InstanceLocation = c.InstanceLocation.AddReferenceToken(name)
				childA, childE := c.Apply()
				annotations = append(annotations, childA...)
				errors = append(errors, childE...)
			}
		}
	}

	// TODO: Add error for patternProperties

	annotations = append(annotations, Annotation{
		InstanceLocation:        ctx.InstanceLocation,
		Keyword:                 ctx.Keyword,
		KeywordLocation:         ctx.KeywordLocation,
		AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
		Value:                   patternPropertiesName,
	})

	return
}
