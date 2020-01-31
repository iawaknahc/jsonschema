package jsonschema

type Required struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
	Missing  []string `json:"missing"`
}

var _ Keyworder = Required{}
var _ Applicator = Required{}

func (_ Required) Keyword() string {
	return "required"
}

func (_ Required) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	obj, ok := ctx.Instance.(map[string]interface{})
	if !ok {
		return
	}

	var expected []string
	for _, name := range UnwrapJSON(ctx.Schema).([]interface{}) {
		expected = append(expected, name.(string))
	}

	actualSet := map[string]struct{}{}
	var actual []string
	for name := range obj {
		actual = append(actual, name)
		actualSet[name] = struct{}{}
	}

	var missing []string
	for _, name := range expected {
		_, ok := actualSet[name]
		if !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		errors = append(errors, Error{
			Keyword:                 ctx.Keyword,
			InstanceLocation:        ctx.InstanceLocation,
			KeywordLocation:         ctx.KeywordLocation,
			AbsoluteKeywordLocation: ctx.AbsoluteKeywordLocation,
			Value: Required{
				Expected: expected,
				Actual:   actual,
				Missing:  missing,
			},
		})
	}

	return
}
