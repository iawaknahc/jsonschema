package jsonschema

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

func MustParseURL(input string) url.URL {
	u, err := url.Parse(input)
	if err != nil {
		panic(err)
	}
	return *u
}

func TestGetSchema(t *testing.T) {
	coll := NewCollection()
	schemaStr := `false`
	coll.AddSchema(strings.NewReader(schemaStr), "")
	schema, err := coll.GetSchema("")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	expected := &JSON{
		JSONValue: false,
	}
	if !reflect.DeepEqual(schema, expected) {
		t.Errorf("%v != %v", schema, expected)
	}
}

func TestAddSchemaAppendixA(t *testing.T) {
	coll := NewCollection()
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
	annotatedSchema := JSON{
		BaseURI:           MustParseURL("https://example.com/root.json"),
		CanonicalLocation: jsonpointer.MustParse(""),
		JSONValue: map[string]JSON{
			"$id": JSON{
				BaseURI:           MustParseURL("https://example.com/root.json"),
				CanonicalLocation: jsonpointer.MustParse("/$id"),
				JSONValue:         "https://example.com/root.json",
			},
			"$defs": JSON{
				BaseURI:           MustParseURL("https://example.com/root.json"),
				CanonicalLocation: jsonpointer.MustParse("/$defs"),
				JSONValue: map[string]JSON{
					"A": JSON{
						BaseURI:           MustParseURL("https://example.com/root.json"),
						CanonicalLocation: jsonpointer.MustParse("/$defs/A"),
						JSONValue: map[string]JSON{
							"$anchor": JSON{
								BaseURI:           MustParseURL("https://example.com/root.json"),
								CanonicalLocation: jsonpointer.MustParse("/$defs/A/$anchor"),
								JSONValue:         "foo",
							},
						}},
					"B": JSON{
						BaseURI:           MustParseURL("https://example.com/other.json"),
						CanonicalLocation: jsonpointer.MustParse(""),
						JSONValue: map[string]JSON{
							"$defs": JSON{
								BaseURI:           MustParseURL("https://example.com/other.json"),
								CanonicalLocation: jsonpointer.MustParse("/$defs"),
								JSONValue: map[string]JSON{
									"X": JSON{
										BaseURI:           MustParseURL("https://example.com/other.json"),
										CanonicalLocation: jsonpointer.MustParse("/$defs/X"),
										JSONValue: map[string]JSON{
											"$anchor": JSON{
												BaseURI:           MustParseURL("https://example.com/other.json"),
												CanonicalLocation: jsonpointer.MustParse("/$defs/X/$anchor"),
												JSONValue:         "bar",
											},
										},
									},
									"Y": JSON{
										BaseURI:           MustParseURL("https://example.com/t/inner.json"),
										CanonicalLocation: jsonpointer.MustParse(""),
										JSONValue: map[string]JSON{
											"$anchor": JSON{
												BaseURI:           MustParseURL("https://example.com/t/inner.json"),
												CanonicalLocation: jsonpointer.MustParse("/$anchor"),
												JSONValue:         "bar",
											},
											"$id": JSON{
												BaseURI:           MustParseURL("https://example.com/t/inner.json"),
												CanonicalLocation: jsonpointer.MustParse("/$id"),
												JSONValue:         "t/inner.json",
											},
										},
									},
								},
							},
							"$id": JSON{
								BaseURI:           MustParseURL("https://example.com/other.json"),
								CanonicalLocation: jsonpointer.MustParse("/$id"),
								JSONValue:         "other.json",
							},
						},
					},
					"C": JSON{
						BaseURI:           MustParseURL("urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f"),
						CanonicalLocation: jsonpointer.MustParse(""),
						JSONValue: map[string]JSON{
							"$id": JSON{
								BaseURI:           MustParseURL("urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f"),
								CanonicalLocation: jsonpointer.MustParse("/$id"),
								JSONValue:         "urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f",
							},
						},
					},
				},
			},
		},
	}
	err := coll.AddSchema(strings.NewReader(schemaStr), "")
	if err != nil {
		t.Errorf("err: %v\n", err)
	}
	cases := []struct {
		ref string
		ptr string
	}{
		{"https://example.com/root.json", "#"},
		{"https://example.com/root.json#", "#"},

		{"https://example.com/root.json#foo", "#/$defs/A"},
		{"https://example.com/root.json#/$defs/A", "#/$defs/A"},

		{"https://example.com/other.json", "#/$defs/B"},
		{"https://example.com/other.json#", "#/$defs/B"},
		{"https://example.com/root.json#/$defs/B", "#/$defs/B"},

		{"https://example.com/other.json#bar", "#/$defs/B/$defs/X"},
		{"https://example.com/other.json#/$defs/X", "#/$defs/B/$defs/X"},
		{"https://example.com/root.json#/$defs/B/$defs/X", "#/$defs/B/$defs/X"},

		{"https://example.com/t/inner.json", "#/$defs/B/$defs/Y"},
		{"https://example.com/t/inner.json#", "#/$defs/B/$defs/Y"},
		{"https://example.com/t/inner.json#bar", "#/$defs/B/$defs/Y"},
		{"https://example.com/other.json#/$defs/Y", "#/$defs/B/$defs/Y"},
		{"https://example.com/root.json#/$defs/B/$defs/Y", "#/$defs/B/$defs/Y"},

		{"urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f", "#/$defs/C"},
		{"urn:uuid:ee564b8a-7a87-4125-8c96-e9f123d6766f#", "#/$defs/C"},
		{"https://example.com/root.json#/$defs/C", "#/$defs/C"},
	}

	for _, c := range cases {
		actual, err := coll.GetSchema(c.ref)
		if err != nil {
			t.Errorf("err: %v\n", err)
		}
		ptr, err := jsonpointer.Parse(c.ptr)
		if err != nil {
			t.Errorf("err: %v\n", err)
		}

		expected, err := TraverseJSON(ptr, annotatedSchema)
		if err != nil {
			t.Errorf("err: %v\n", err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%s != %s\n", c.ref, c.ptr)
			t.Errorf("actual: %#v\n", actual)
			t.Errorf("expected: %#v\n", expected)
		}
	}
}
