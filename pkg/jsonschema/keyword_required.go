package jsonschema

import (
	"sort"
)

type Required struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
	Missing  []string `json:"missing"`
}

var _ Keyword = Required{}

func (_ Required) Keyword() string {
	return "required"
}

func (_ Required) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	var expected []string
	for _, name := range UnwrapJSON(input.Schema).([]interface{}) {
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

	// Sort them to ensure the order is stable.
	// It is very useful if this node is testing data.
	sort.Strings(expected)
	sort.Strings(actual)
	sort.Strings(missing)

	if len(missing) > 0 {
		input.Valid = false
		input.Info = Required{
			Expected: expected,
			Actual:   actual,
			Missing:  missing,
		}
	}

	return &input, nil
}
