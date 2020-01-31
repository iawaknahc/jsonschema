package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

func TestWrapJSON(t *testing.T) {
	schemaStr := `
	{
		"$id": "https://example.com/root.json",
		"$defs": {
			"A": { "$anchor": "foo" },
			"B": {
				"$id": "other.json",
				"$defs": {
					"X": { "$anchor": "bar" },
					"Y": {
						"$id": "t/inner.json",
						"$anchor": "bar"
					}
				}
			},
			"C": {
				"$id": "urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f"
			}
		}
	}
	`
	var schema interface{}
	err := json.Unmarshal([]byte(schemaStr), &schema)
	if err != nil {
		t.Errorf("expected err: %v", err)
	}

	actual := WrapJSON(schema)
	expected := JSON{
		JSONValue: map[string]JSON{
			"$id": JSON{JSONValue: "https://example.com/root.json"},
			"$defs": JSON{
				JSONValue: map[string]JSON{
					"A": JSON{
						JSONValue: map[string]JSON{
							"$anchor": JSON{JSONValue: "foo"},
						}},
					"B": JSON{
						JSONValue: map[string]JSON{
							"$defs": JSON{
								JSONValue: map[string]JSON{
									"X": JSON{
										JSONValue: map[string]JSON{
											"$anchor": JSON{JSONValue: "bar"},
										},
									},
									"Y": JSON{
										JSONValue: map[string]JSON{
											"$anchor": JSON{JSONValue: "bar"},
											"$id":     JSON{JSONValue: "t/inner.json"},
										},
									},
								},
							},
							"$id": JSON{JSONValue: "other.json"},
						},
					},
					"C": JSON{
						JSONValue: map[string]JSON{
							"$id": JSON{JSONValue: "urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f"},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual != expected")
	}

	unwrapped := UnwrapJSON(actual)
	if !reflect.DeepEqual(unwrapped, schema) {
		t.Errorf("unwrapped != schema")
	}
}

func TestTraverseJSON(t *testing.T) {
	schemaStr := `
	{
		"$id": "https://example.com/root.json",
		"$defs": {
			"A": { "$anchor": "foo" },
			"B": {
				"$id": "other.json",
				"$defs": {
					"X": { "$anchor": "bar" },
					"Y": {
						"$id": "t/inner.json",
						"$anchor": "bar"
					}
				}
			},
			"C": {
				"$id": "urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f"
			}
		}
	}
	`
	var schema interface{}
	err := json.Unmarshal([]byte(schemaStr), &schema)
	if err != nil {
		t.Errorf("expected err: %v", err)
	}

	j := WrapJSON(schema)
	ptr, err := jsonpointer.Parse("/$defs/B/$id")
	if err != nil {
		t.Errorf("expected err: %v", err)
	}

	actual, err := TraverseJSON(ptr, j)
	if err != nil {
		t.Errorf("expected err: %v", err)
	}

	expected := &JSON{
		JSONValue: "other.json",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%#v != %#v", actual, expected)
	}
}
