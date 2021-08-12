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
		collection := NewMetaschemaCollection()
		err = collection.AddSchema(bytes.NewReader(schemaBytes), "")
		if err != nil {
			t.Fatalf("failed to add schema: %v", err)
		}

		for _, test := range c.Tests {
			shouldSkip := false
			for _, s := range skip {
				if test.Description == s {
					shouldSkip = true
					break
				}
			}
			if shouldSkip {
				continue
			}

			dataBytes, err := json.Marshal(test.Data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			node, err := collection.Apply(id, bytes.NewReader(dataBytes))
			if err != nil {
				if test.Valid {
					t.Errorf("%s: %s: treating error as invalid: %v", c.Description, test.Description, err)
				}
			} else {
				if test.Valid && !node.Valid {
					t.Fatalf("%s: %s: %+v", c.Description, test.Description, node)
				}
				if !test.Valid && node.Valid {
					t.Fatalf("%s: %s", c.Description, test.Description)
				}
			}
		}
	}
}

func TestBoolean(t *testing.T) {
	test(t, "draft2019-09/boolean_schema.json")
}

func TestRef(t *testing.T) {
	// FIXME: https://github.com/json-schema-org/JSON-Schema-Test-Suite/pull/327
	test(t, "draft2019-09/ref.json",
		"ref creates new scope when adjacent to keywords",
	)
}

// FIXME: Use refRemote.json
// func TestRefRemote(t *testing.T) {
// 	test(t, "draft2019-09/refRemote.json")
// }

func TestID(t *testing.T) {
	test(t, "draft2019-09/id.json",
		// FIXME: schema collection
		"$id inside an enum is not a real identifier",
	)
}

func TestAnchor(t *testing.T) {
	test(t, "draft2019-09/anchor.json",
		// FIXME: schema collection
		"$anchor inside an enum is not a real identifier",
	)
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

func TestMinContains(t *testing.T) {
	test(t, "draft2019-09/minContains.json")
}

func TestMaxContains(t *testing.T) {
	test(t, "draft2019-09/maxContains.json")
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
	test(t, "draft2019-09/format.json",
		"invalid ipv4 string is only an annotation by default",
	)
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

func TestDefs(t *testing.T) {
	test(t, "draft2019-09/defs.json")
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

// FIXME: content.json
// func TestContent(t *testing.T) {
// 	test(t, "draft2019-09/optional/content.json")
// }

// FIXME: ecmascript-regex.json
// func TestECMAScriptRegex(t *testing.T) {
// 	test(t, "draft2019-09/optional/ecmascript-regex.json")
// }

// FIXME: multipleOf
// func TestFloatOverflow(t *testing.T) {
// 	test(t, "draft2019-09/optional/float-overflow.json")
// }

func TestNonBMPRegex(t *testing.T) {
	test(t, "draft2019-09/optional/non-bmp-regex.json")
}

func TestFormatIPv4(t *testing.T) {
	test(t, "draft2019-09/optional/format/ipv4.json")
}

func TestInfiniteLoopDetection(t *testing.T) {
	test(t, "draft2019-09/infinite-loop-detection.json")
}

func TestUnknownKeyword(t *testing.T) {
	test(t, "draft2019-09/unknownKeyword.json",
		// FIXME: schema collection
		"$id inside an unknown keyword is not a real identifier",
	)
}

func TestRecursiveRef(t *testing.T) {
	test(t, "draft2019-09/recursiveRef.json",
		// FIXME: recursiveRef
		"multiple dynamic paths to the $recursiveRef keyword",
		"dynamic $recursiveRef destination (not predictable at schema compile time)",
	)
}
