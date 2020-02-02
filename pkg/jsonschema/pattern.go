package jsonschema

type Pattern struct {
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

var _ Keyword = Pattern{}

func (_ Pattern) Keyword() string {
	return "pattern"
}

func (_ Pattern) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	str, ok := input.Instance.(string)
	if !ok {
		return &input, nil
	}
	pattern := input.Schema.JSONValue.(string)
	re, err := ctx.CompilePattern(pattern)
	if err != nil {
		return nil, err
	}
	if !re.MatchString(str) {
		input.Valid = false
		input.Info = Pattern{
			Expected: pattern,
			Actual:   str,
		}
	}
	return &input, nil
}
