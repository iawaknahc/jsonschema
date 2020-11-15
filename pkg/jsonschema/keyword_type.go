package jsonschema

import (
	"encoding/json"
	"sort"

	"github.com/cockroachdb/apd/v2"
)

type Type struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
}

var _ Keyword = Type{}

func (_ Type) Keyword() string {
	return "type"
}

func (_ Type) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	// Prepare expected.
	var expected []string
	switch a := UnwrapJSON(input.Scope.Schema).(type) {
	case string:
		expected = []string{a}
	case []interface{}:
		for _, b := range a {
			if c, ok := b.(string); ok {
				expected = append(expected, c)
			}
		}
	}

	// Prepare actual
	var actual []string
	switch d := input.Instance.(type) {
	case nil:
		actual = append(actual, "null")
	case string:
		actual = append(actual, "string")
	case bool:
		actual = append(actual, "boolean")
	case map[string]interface{}:
		actual = append(actual, "object")
	case []interface{}:
		actual = append(actual, "array")
	case json.Number:
		decimal, _, err := apd.NewFromString(string(d))
		if err != nil {
			return nil, err
		}
		res, err := apd.BaseContext.RoundToIntegralExact(decimal, decimal)
		if err != nil {
			return nil, err
		}
		integer := !res.Inexact()
		if integer {
			actual = append(actual, "number", "integer")
		} else {
			actual = append(actual, "number")
		}
	}

	// Sort them to ensure the order is stable.
	// It is very useful if this node is testing data.
	sort.Strings(expected)
	sort.Strings(actual)

	intersection := intersectString(expected, actual)
	if len(intersection) <= 0 {
		input.Valid = false
		input.Info = Type{
			Expected: expected,
			Actual:   actual,
		}
	}

	return &input, nil
}

func intersectString(a []string, b []string) (out []string) {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				out = append(out, x)
			}
		}
	}
	return
}
