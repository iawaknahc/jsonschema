# Caveat

The regular expression syntax is not EMCA 262 but [RE2](https://golang.org/s/re2syntax).
This effects the keywords `patternProperties` and `pattern`.

# Output

- [x] Flag
- [ ] Basic
- [ ] Detailed
- [x] Verbose

# Keyword Support

- [ ] $schema
- [ ] $vocabulary
- [x] $id
- [x] $anchor
- [x] $ref
- [x] $defs
- [x] $comment

- [ ] $recursiveRef
- [ ] $recursiveAnchor

- [x] allOf
- [x] anyOf
- [x] oneOf
- [x] not

- [x] if
- [x] then
- [x] else
- [ ] dependentSchemas (no test in official test suite yet)

- [x] items
- [x] additionalItems
- [ ] unevaluatedItems (no test in official test suite yet)
- [ ] contains

- [x] properties
- [x] patternProperties
- [x] additionalProperties
- [ ] unevaluatedProperties
- [x] propertyNames

- [x] type
- [x] enum
- [x] const
- [x] multipleOf
- [x] maximum
- [x] exclusiveMaximum
- [x] minimum
- [x] exclusiveMinimum
- [x] maxLength
- [x] minLength
- [x] pattern
- [x] maxItems
- [x] minItems
- [ ] uniqueItems
- [ ] maxContains (no test in official test suite yet)
- [ ] minContains (no test in official test suite yet)
- [x] maxProperties
- [x] minProperties
- [x] required
- [ ] dependentRequired (no test in official test suite yet)

- [ ] format
- [ ] title
- [ ] description
- [ ] default
- [ ] deprecated
- [ ] readOnly
- [ ] writeOnly
- [ ] examples

- [ ] contentEncoding
- [ ] contentMediaType
- [ ] contentSchema
