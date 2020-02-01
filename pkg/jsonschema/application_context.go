package jsonschema

import (
	"strings"
)

type ErrCircularReference struct {
	Locations []Location
}

func (e ErrCircularReference) Error() string {
	strs := make([]string, len(e.Locations))
	for i, l := range e.Locations {
		strs[i] = l.String()
	}
	return strings.Join(strs, " -> ")
}

type ApplicationContext struct {
	Collection         *Collection
	Vocabulary         Vocabulary
	ReferencedLocation []Location
}

func (c ApplicationContext) Apply(input Node) (*Node, error) {
	// Handle $ref
	if obj, ok := input.Schema.JSONValue.(map[string]JSON); ok {
		if ref, ok := obj["$ref"].JSONValue.(string); ok {
			u, err := input.Schema.BaseURI.Parse(ref)
			if err != nil {
				return nil, err
			}
			referencedSchema, err := c.Collection.GetSchema(u.String())
			if err != nil {
				return nil, err
			}
			location := Location{
				BaseURI:     referencedSchema.BaseURI,
				JSONPointer: referencedSchema.CanonicalLocation,
			}
			childInput := Node{
				Valid:                   true,
				Parent:                  &input,
				Instance:                input.Instance,
				InstanceLocation:        input.InstanceLocation,
				Schema:                  *referencedSchema,
				Keyword:                 "$ref",
				KeywordLocation:         input.KeywordLocation.AddReferenceToken("$ref"),
				AbsoluteKeywordLocation: location,
			}

			// Detect cycle.
			for _, l := range c.ReferencedLocation {
				if l.String() == location.String() {
					return nil, ErrCircularReference{
						Locations: c.ReferencedLocation,
					}
				}
			}
			c.ReferencedLocation = append(c.ReferencedLocation, location)

			child, err := c.Apply(childInput)
			if err != nil {
				return nil, err
			}

			input.Valid = child.Valid
			input.Children = append(input.Children, *child)

			return &input, nil
		}
	}
	c.ReferencedLocation = nil

	// Handle boolean schema
	if b, ok := input.Schema.JSONValue.(bool); ok {
		// TODO(boolean): fill in info
		input.Valid = b
		return &input, nil
	}

	// Handle each keywords
	if schema, ok := input.Schema.JSONValue.(map[string]JSON); ok {
		// We need to apply the keywords with the order in the vocabulary.
		// We also need to ignore any unknown keywords.
		keywords := map[string]struct{}{}
		for name := range schema {
			keywords[name] = struct{}{}
		}
		// We now have a set of present keywords in the schema object.
		// Process them in the vocabulary order.
		for _, keyword := range c.Vocabulary.Keywords {
			k := keyword.Keyword()
			// keyword not found in this schema object.
			// Skip to the next keyword.
			if _, ok := keywords[k]; !ok {
				continue
			}
			// Remove processed keywords.
			delete(keywords, k)
			childInput := Node{
				Valid:                   true,
				Parent:                  &input,
				Instance:                input.Instance,
				InstanceLocation:        input.InstanceLocation,
				Schema:                  schema[k],
				Keyword:                 k,
				KeywordLocation:         input.KeywordLocation.AddReferenceToken(k),
				AbsoluteKeywordLocation: input.AbsoluteKeywordLocation.AddReferenceToken(k),
			}
			child, err := keyword.Apply(c, childInput)
			if err != nil {
				return nil, err
			}

			if !child.Valid {
				input.Valid = false
			}
			input.Children = append(input.Children, *child)
		}
		// We now have a set of unknown keywords in the schema object.
		// Ignore them by assuming valid.
		for keyword := range keywords {
			child := Node{
				Valid:                   true,
				Parent:                  &input,
				Instance:                input.Instance,
				InstanceLocation:        input.InstanceLocation,
				Schema:                  schema[keyword],
				Keyword:                 keyword,
				KeywordLocation:         input.KeywordLocation.AddReferenceToken(keyword),
				AbsoluteKeywordLocation: input.AbsoluteKeywordLocation.AddReferenceToken(keyword),
			}
			input.Children = append(input.Children, child)
		}

		return &input, nil
	}

	// The schema is neither boolean nor object.
	return nil, ErrNotASchema
}
