package jsonschema

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type Case struct {
	Description string      `json:"description"`
	Schema      interface{} `json:"schema"`
	Tests       []Test      `json:"tests"`
}

type Test struct {
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	Valid       bool        `json:"valid"`
}

func test(t *testing.T, p string, skip ...string) {
	p = filepath.Join("../../tests", p)
	f, err := os.Open(p)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	var cases []Case
	d := json.NewDecoder(f)
	d.UseNumber()
	err = d.Decode(&cases)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	for _, c := range cases {
		shouldSkip := false
		for _, s := range skip {
			if c.Description == s {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}

		// Special handling for ref.json
		id := ""
		if obj, ok := c.Schema.(map[string]interface{}); ok {
			if i, ok := obj["$id"].(string); ok {
				id = i
			}
		}

		schemaBytes, err := json.Marshal(c.Schema)
		if err != nil {
			t.Fatalf("failed to marshal schema: %v", err)
		}
		collection := NewCollection()
		err = collection.AddSchema(bytes.NewReader(schemaBytes), "")
		if err != nil {
			t.Fatalf("failed to add schema: %v", err)
		}

		for _, test := range c.Tests {
			dataBytes, err := json.Marshal(test.Data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			node, err := collection.Apply(id, bytes.NewReader(dataBytes))
			if err != nil {
				t.Fatalf("%s: %s: unexpected error: %v", c.Description, test.Description, err)
			}
			if test.Valid && !node.Valid {
				t.Fatalf("%s: %s: %+v", c.Description, test.Description, node)
			}
			if !test.Valid && node.Valid {
				t.Fatalf("%s: %s", c.Description, test.Description)
			}
		}
	}
}

func TestBoolean(t *testing.T) {
	test(t, "draft2019-09/boolean_schema.json")
}

func TestRef(t *testing.T) {
	// TODO: embed the meta-schema of draft 2019-09
	test(t, "draft2019-09/ref.json", "remote ref, containing refs itself")
}

func TestAnchor(t *testing.T) {
	test(t, "draft2019-09/anchor.json")
}

func TestType(t *testing.T) {
	test(t, "draft2019-09/type.json")
}

func TestConst(t *testing.T) {
	test(t, "draft2019-09/const.json")
}

func TestEnum(t *testing.T) {
	test(t, "draft2019-09/enum.json")
}

func TestMaxItems(t *testing.T) {
	test(t, "draft2019-09/maxItems.json")
}

func TestMinItems(t *testing.T) {
	test(t, "draft2019-09/minItems.json")
}

func TestContains(t *testing.T) {
	test(t, "draft2019-09/contains.json")
}

func TestUniqueItems(t *testing.T) {
	test(t, "draft2019-09/uniqueItems.json")
}

func TestMultipleOf(t *testing.T) {
	test(t, "draft2019-09/multipleOf.json")
}

func TestMaximum(t *testing.T) {
	test(t, "draft2019-09/maximum.json")
}

func TestExclusiveMaximum(t *testing.T) {
	test(t, "draft2019-09/exclusiveMaximum.json")
}

func TestMinimum(t *testing.T) {
	test(t, "draft2019-09/minimum.json")
}

func TestExclusiveMinimum(t *testing.T) {
	test(t, "draft2019-09/exclusiveMinimum.json")
}

func TestPropertyNames(t *testing.T) {
	test(t, "draft2019-09/propertyNames.json")
}

func TestMaxProperties(t *testing.T) {
	test(t, "draft2019-09/maxProperties.json")
}

func TestMinProperties(t *testing.T) {
	test(t, "draft2019-09/minProperties.json")
}

func TestItems(t *testing.T) {
	test(t, "draft2019-09/items.json")
}

func TestAdditionalItems(t *testing.T) {
	test(t, "draft2019-09/additionalItems.json")
}

func TestProperties(t *testing.T) {
	test(t, "draft2019-09/properties.json")
}

func TestPatternProperties(t *testing.T) {
	test(t, "draft2019-09/patternProperties.json")
}

func TestAdditionalProperties(t *testing.T) {
	test(t, "draft2019-09/additionalProperties.json")
}

func TestRequired(t *testing.T) {
	test(t, "draft2019-09/required.json")
}

func TestDependentRequired(t *testing.T) {
	test(t, "draft2019-09/dependentRequired.json")
}

func TestMaxLength(t *testing.T) {
	test(t, "draft2019-09/maxLength.json")
}

func TestMinLength(t *testing.T) {
	test(t, "draft2019-09/minLength.json")
}

func TestPattern(t *testing.T) {
	test(t, "draft2019-09/pattern.json")
}

func TestFormat(t *testing.T) {
	test(t, "draft2019-09/format.json")
}

func TestFormatIPv4(t *testing.T) {
	test(t, "draft2019-09/optional/format/ipv4.json")
}

func TestAllOf(t *testing.T) {
	test(t, "draft2019-09/allOf.json")
}

func TestOneOf(t *testing.T) {
	test(t, "draft2019-09/oneOf.json")
}

func TestAnyOf(t *testing.T) {
	test(t, "draft2019-09/anyOf.json")
}

func TestNot(t *testing.T) {
	test(t, "draft2019-09/not.json")
}

func TestIfThenElse(t *testing.T) {
	test(t, "draft2019-09/if-then-else.json")
}

func TestDependentSchemas(t *testing.T) {
	test(t, "draft2019-09/dependentSchemas.json")
}

func TestDefault(t *testing.T) {
	test(t, "draft2019-09/default.json")
}

func TestBignum(t *testing.T) {
	test(t, "draft2019-09/optional/bignum.json")
}

func TestZeroTerminatedFloats(t *testing.T) {
	test(t, "draft2019-09/optional/zeroTerminatedFloats.json")
}

func TestRefOfUnknownKeyword(t *testing.T) {
	test(t, "draft2019-09/optional/refOfUnknownKeyword.json")
}
