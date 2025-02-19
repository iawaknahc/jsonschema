package jsonschema

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func Print(t *testing.T, msg string, actual OutputNode, expected OutputNode) {
	t.Errorf("%s", msg)
	b, _ := json.MarshalIndent(actual, "", "  ")
	t.Errorf("actual: %s", string(b))
	b, _ = json.MarshalIndent(expected, "", "  ")
	t.Errorf("expected: %s", string(b))
}

func OutputNodeEqual(t *testing.T, actual OutputNode, expected OutputNode) bool {
	if actual.Valid != expected.Valid {
		Print(t, "valid", actual, expected)
		return false
	}
	if actual.KeywordLocation != expected.KeywordLocation {
		Print(t, "keywordLocation", actual, expected)
		return false
	}
	if actual.InstanceLocation != expected.InstanceLocation {
		Print(t, "instanceLocation", actual, expected)
		return false
	}
	if expected.AbsoluteKeywordLocation != "" && actual.AbsoluteKeywordLocation != expected.AbsoluteKeywordLocation {
		Print(t, "absoluteKeywordLocation", actual, expected)
		return false
	}

	numErrors := 0
	for _, a := range actual.Errors {
		keywordLocation := a.KeywordLocation
		instanceLocation := a.InstanceLocation
		for _, e := range expected.Errors {
			if keywordLocation == e.KeywordLocation && instanceLocation == e.InstanceLocation {
				numErrors++
				if !OutputNodeEqual(t, a, e) {
					return false
				}
			}
		}
	}
	if numErrors < len(expected.Errors) {
		Print(t, "errors", actual, expected)
		return false
	}

	numAnnotations := 0
	for _, a := range actual.Annotations {
		keywordLocation := a.KeywordLocation
		instanceLocation := a.InstanceLocation
		for _, e := range expected.Annotations {
			if keywordLocation == e.KeywordLocation && instanceLocation == e.InstanceLocation {
				numAnnotations++
				if !OutputNodeEqual(t, a, e) {
					return false
				}
			}
		}
	}
	if numAnnotations < len(expected.Annotations) {
		Print(t, "annotations", actual, expected)
		return false
	}

	return true
}

func TestVerboseSimple(t *testing.T) {
	schema := `
	{
		"$id": "https://example.com/polygon",
		"$schema": "https://json-schema.org/draft/2019-09/schema",
		"type": "object",
		"properties": {
			"validProp": true
		},
		"additionalProperties": false
	}
	`
	instance := `
	{
		"validProp": 5,
		"disallowedProp": "value"
	}
	`
	expectedFile := "output_verbose_simple.json"

	var err error
	var expected OutputNode
	f, err := os.Open(expectedFile)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	err = json.NewDecoder(f).Decode(&expected)
	if err != nil {
		t.Fatalf("failed to unmarshal expected: %v", err)
	}

	collection := NewCollection()
	err = collection.AddSchema(strings.NewReader(schema), "")
	if err != nil {
		t.Fatalf("failed to add schema: %v", err)
	}
	ctx := context.Background()
	node, err := collection.Apply(ctx, "https://example.com/polygon", strings.NewReader(instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	actual := node.Verbose()

	OutputNodeEqual(t, actual, expected)
}

func TestVerbose(t *testing.T) {
	schema := `
	{
		"$id": "https://example.com/polygon",
		"$schema": "https://json-schema.org/draft/2019-09/schema",
		"$defs": {
			"point": {
				"type": "object",
				"properties": {
					"x": { "type": "number" },
					"y": { "type": "number" }
				},
				"additionalProperties": false,
				"required": [ "x", "y" ]
			}
		},
		"type": "array",
		"items": { "$ref": "#/$defs/point" },
		"minItems": 3
	}
	`
	instance := `
	[
	{
		"x": 2.5,
		"y": 1.3
	},
	{
		"x": 1,
		"z": 6.7
	}
	]
	`
	expectedFile := "output_verbose.json"

	var err error
	var expected OutputNode
	f, err := os.Open(expectedFile)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	err = json.NewDecoder(f).Decode(&expected)
	if err != nil {
		t.Fatalf("failed to unmarshal expected: %v", err)
	}

	collection := NewCollection()
	err = collection.AddSchema(strings.NewReader(schema), "")
	if err != nil {
		t.Fatalf("failed to add schema: %v", err)
	}
	ctx := context.Background()
	node, err := collection.Apply(ctx, "https://example.com/polygon", strings.NewReader(instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	actual := node.Verbose()

	OutputNodeEqual(t, actual, expected)
}
