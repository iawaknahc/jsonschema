# Caveat

The regular expression syntax is not EMCA 262 but [RE2](https://golang.org/s/re2syntax).
This effects the keywords `patternProperties` and `pattern`.

# Output

- [x] Flag
- [ ] Basic
- [ ] Detailed
- [x] Verbose

# Keyword Support

- [x] $schema
- [ ] $vocabulary
- [x] $id
- [x] $anchor
- [x] $ref
- [x] $defs
- [x] $comment
- [x] $recursiveRef
- [x] $recursiveAnchor
- [x] allOf
- [x] anyOf
- [x] oneOf
- [x] not
- [x] if
- [x] then
- [x] else
- [x] dependentSchemas
- [x] items
- [x] additionalItems
- [ ] unevaluatedItems
- [x] contains
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
- [x] uniqueItems
- [x] maxContains
- [x] minContains
- [x] maxProperties
- [x] minProperties
- [x] required
- [x] dependentRequired
- [x] format
  - [ ] date-time
  - [ ] date
  - [ ] time
  - [ ] duration
  - [x] email
  - [ ] idn-email
  - [ ] hostname
  - [ ] idn-hostname
  - [x] ipv4
  - [ ] ipv6
  - [ ] uri
  - [ ] uri-reference
  - [ ] iri
  - [ ] iri-reference
  - [ ] uuid
  - [ ] uri-template
  - [ ] json-pointer
  - [ ] relative-json-pointer
  - [ ] regex
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
