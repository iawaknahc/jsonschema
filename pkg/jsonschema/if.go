package jsonschema

type If struct{}

var _ Keyword = If{}
var _ AnnotatingKeyword = If{}

func (_ If) Keyword() string {
	return "if"
}

func (_ If) CombineAnnotations(values []interface{}) (interface{}, bool) {
	if len(values) <= 0 {
		return nil, false
	}
	if len(values) == 1 {
		return values[0], true
	}
	panic("impossible to combine annotation value of if")
}

func (_ If) Apply(ctx ApplicationContext, input Node) (*Node, error) {
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
	// if is always valid
	input.Valid = true
	input.Children = append(input.Children, *child)
	input.Annotation = child.Valid
	return &input, nil
}
