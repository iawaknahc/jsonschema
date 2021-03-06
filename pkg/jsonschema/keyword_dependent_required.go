package jsonschema

import (
	"sort"
)

type DependentRequired struct {
	Required map[string]Required
}

var _ Keyword = DependentRequired{}

func (_ DependentRequired) Keyword() string {
	return "dependentRequired"
}

func (_ DependentRequired) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	obj, ok := input.Instance.(map[string]interface{})
	if !ok {
		return &input, nil
	}

	required := map[string]Required{}

	actualSet := map[string]struct{}{}
	var actual []string
	for name := range obj {
		actual = append(actual, name)
		actualSet[name] = struct{}{}
	}

	// Sort it to ensure the order is stable.
	// It is very useful if this node is testing data.
	sort.Strings(actual)

	for name, schema := range input.Scope.Schema.JSONValue.(map[string]JSON) {
		_, ok := obj[name]
		if !ok {
			continue
		}

		var expected []string
		for _, requiredName := range UnwrapJSON(schema).([]interface{}) {
			expected = append(expected, requiredName.(string))
		}

		var missing []string
		for _, requiredName := range expected {
			_, ok := actualSet[requiredName]
			if !ok {
				missing = append(missing, requiredName)
			}
		}

		// Sort them to ensure the order is stable.
		// It is very useful if this node is testing data.
		sort.Strings(expected)
		sort.Strings(missing)

		if len(missing) > 0 {
			required[name] = Required{
				Expected: expected,
				Actual:   actual,
				Missing:  missing,
			}
		}
	}

	if len(required) > 0 {
		input.Valid = false
		input.Info = DependentRequired{
			Required: required,
		}
	}

	return &input, nil
}
