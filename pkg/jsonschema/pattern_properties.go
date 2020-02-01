package jsonschema

import (
	"regexp"
)

type PatternProperties struct{}

var _ Keyword = PatternProperties{}
var _ AnnotatingKeyword = PatternProperties{}

func (_ PatternProperties) Keyword() string {
	return "patternProperties"
}

func (_ PatternProperties) CombineAnnotations(values []interface{}) (interface{}, bool) {
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

func (_ PatternProperties) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	patternPropertiesName := map[string]struct{}{}
	for pattern, schema := range input.Schema.JSONValue.(map[string]JSON) {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		for name, val := range obj {
			if re.MatchString(name) {
				patternPropertiesName[name] = struct{}{}
				childInput := Node{
					Valid:                   true,
					Parent:                  &input,
					Instance:                val,
					InstanceLocation:        input.InstanceLocation.AddReferenceToken(name),
					Schema:                  schema,
					KeywordLocation:         input.KeywordLocation.AddReferenceToken(pattern),
					AbsoluteKeywordLocation: input.AbsoluteKeywordLocation.AddReferenceToken(pattern),
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
	}

	input.Annotation = patternPropertiesName

	return &input, nil
}
