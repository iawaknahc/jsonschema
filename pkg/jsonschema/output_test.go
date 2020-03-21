package jsonschema

import (
	"encoding/json"
	"strings"
	"testing"
)

func OutputNodeEqual(actual OutputNode, expected OutputNode) bool {
	if actual.Valid != expected.Valid {
		return false
	}
	if actual.KeywordLocation != expected.KeywordLocation {
		return false
	}
	if actual.InstanceLocation != expected.InstanceLocation {
		return false
	}
	if expected.AbsoluteKeywordLocation != "" && actual.AbsoluteKeywordLocation != expected.AbsoluteKeywordLocation {
		return false
	}

	numErrors := 0
	for _, a := range actual.Errors {
		keywordLocation := a.KeywordLocation
		instanceLocation := a.InstanceLocation
		for _, e := range expected.Errors {
			if keywordLocation == e.KeywordLocation && instanceLocation == e.InstanceLocation {
				numErrors++
				if !OutputNodeEqual(a, e) {
					return false
				}
			}
		}
	}
	if numErrors < len(expected.Errors) {
		return false
	}

	numAnnotations := 0
	for _, a := range actual.Annotations {
		keywordLocation := a.KeywordLocation
		instanceLocation := a.InstanceLocation
		for _, e := range expected.Annotations {
			if keywordLocation == e.KeywordLocation && instanceLocation == e.InstanceLocation {
				numAnnotations++
				if !OutputNodeEqual(a, e) {
					return false
				}
			}
		}
	}
	if numAnnotations < len(expected.Annotations) {
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
	expectedStr := `
	{
		"valid": false,
		"keywordLocation": "#",
		"instanceLocation": "#",
		"errors": [
		{
			"valid": true,
			"keywordLocation": "#/type",
			"instanceLocation": "#"
		},
		{
			"valid": true,
			"keywordLocation": "#/properties",
			"instanceLocation": "#"
		},
		{
			"valid": false,
			"keywordLocation": "#/additionalProperties",
			"instanceLocation": "#",
			"errors": [
			{
				"valid": false,
				"keywordLocation": "#/additionalProperties",
				"instanceLocation": "#/disallowedProp",
				"error": "Additional property 'disallowedProp' found but was invalid."
			}
			]
		}
		]
	}
	`

	var err error
	var expected OutputNode
	err = json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Fatalf("failed to unmarshal expected: %v", err)
	}

	collection := NewCollection()
	err = collection.AddSchema(strings.NewReader(schema), "")
	if err != nil {
		t.Fatalf("failed to add schema: %v", err)
	}
	node, err := collection.Apply("https://example.com/polygon", strings.NewReader(instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	actual := node.Verbose()

	if !OutputNodeEqual(actual, expected) {
		t.Fatalf("actual != expected")
	}
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
	expectedStr := `
	{
		"valid": false,
		"keywordLocation": "#",
		"instanceLocation": "#",
		"errors": [
		{
			"valid": true,
			"keywordLocation": "#/$defs",
			"instanceLocation": "#"
		},
		{
			"valid": true,
			"keywordLocation": "#/type",
			"instanceLocation": "#"
		},
		{
			"valid": false,
			"keywordLocation": "#/items",
			"instanceLocation": "#",
			"errors": [
			{
				"valid": true,
				"keywordLocation": "#/items/$ref",
				"absoluteKeywordLocation": "https://example.com/polygon#/items/$ref",
				"instanceLocation": "#/0",
				"annotations": [
				{
					"valid": true,
					"keywordLocation": "#/items/$ref",
					"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point",
					"instanceLocation": "#/0",
					"annotations": [
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/type",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/type",
						"instanceLocation": "#/0"
					},
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/properties",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/properties",
						"instanceLocation": "#/0"
					},
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/required",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/required",
						"instanceLocation": "#/0"
					},
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/additionalProperties",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/additionalProperties",
						"instanceLocation": "#/0"
					}
					]
				}
				]
			},
			{
				"valid": false,
				"keywordLocation": "#/items/$ref",
				"absoluteKeywordLocation": "https://example.com/polygon#/items/$ref",
				"instanceLocation": "#/1",
				"errors": [
				{
					"valid": false,
					"keywordLocation": "#/items/$ref",
					"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point",
					"instanceLocation": "#/1",
					"errors": [
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/type",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/type",
						"instanceLocation": "#/1"
					},
					{
						"valid": true,
						"keywordLocation": "#/items/$ref/properties",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/properties",
						"instanceLocation": "#/1"
					},
					{
						"valid": false,
						"keywordLocation": "#/items/$ref/required",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/required",
						"instanceLocation": "#/1"
					},
					{
						"valid": false,
						"keywordLocation": "#/items/$ref/additionalProperties",
						"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/additionalProperties",
						"instanceLocation": "#/1",
						"errors": [
						{
							"valid": false,
							"keywordLocation": "#/items/$ref/additionalProperties",
							"absoluteKeywordLocation": "https://example.com/polygon#/$defs/point/additionalProperties",
							"instanceLocation": "#/1/z"
						}
						]
					}
					]
				}
				]
			}
			]
		},
		{
			"valid": false,
			"keywordLocation": "#/minItems",
			"instanceLocation": "#"
		}
		]
	}
	`

	var err error
	var expected OutputNode
	err = json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Fatalf("failed to unmarshal expected: %v", err)
	}

	collection := NewCollection()
	err = collection.AddSchema(strings.NewReader(schema), "")
	if err != nil {
		t.Fatalf("failed to add schema: %v", err)
	}
	node, err := collection.Apply("https://example.com/polygon", strings.NewReader(instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	actual := node.Verbose()

	if !OutputNodeEqual(actual, expected) {
		t.Fatalf("actual != expected")
	}
}
