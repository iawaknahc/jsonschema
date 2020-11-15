package jsonschema

import (
	"regexp"
	"strings"
	"sync"
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
	Collection   *Collection
	Vocabulary   Vocabulary
	PatternCache *sync.Map
}

func (c ApplicationContext) CompilePattern(pattern string) (*regexp.Regexp, error) {
	value, ok := c.PatternCache.Load(pattern)
	if ok {
		switch v := value.(type) {
		case *regexp.Regexp:
			return v, nil
		case error:
			return nil, v
		default:
			panic("unreachable")
		}
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		c.PatternCache.Store(pattern, err)
	} else {
		c.PatternCache.Store(pattern, re)
	}
	return re, err
}

var handledKeywords map[string]struct{} = map[string]struct{}{
	"$id":     {},
	"$anchor": {},
	"$schema": {},
	// TODO: Handle $vocabulary
	"$vocabulary":      {},
	"$recursiveAnchor": {},
	"minContains":      {},
	"maxContains":      {},
}

func (c ApplicationContext) Apply(input Node) (*Node, error) {
	// Handle boolean schema
	if b, ok := input.Scope.Schema.JSONValue.(bool); ok {
		input.Valid = b
		return &input, nil
	}

	// Handle each keywords
	if schema, ok := input.Scope.Schema.JSONValue.(map[string]JSON); ok {
		// We need to apply the keywords with the order in the vocabulary.
		// We also need to ignore any unknown keywords.
		keywords := map[string]struct{}{}
		for name := range schema {
			// Ignore handled keyword
			if _, ok := handledKeywords[name]; ok {
				continue
			}
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
				Valid:            true,
				Parent:           &input,
				Instance:         input.Instance,
				InstanceLocation: input.InstanceLocation,
				Keyword:          k,
				Scope:            input.Scope.Descend(k, schema[k]),
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
				Valid:            true,
				Parent:           &input,
				Instance:         input.Instance,
				InstanceLocation: input.InstanceLocation,
				Keyword:          keyword,
				Scope:            input.Scope.Descend(keyword, schema[keyword]),
			}
			input.Children = append(input.Children, child)
		}

		return &input, nil
	}

	// The schema is neither boolean nor object.
	return nil, ErrNotASchema
}
