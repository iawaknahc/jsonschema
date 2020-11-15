package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

// Scope stores the necessary information for representing https://json-schema.org/draft/2019-09/json-schema-core.html#scopes
// Scope should be immutable.
type Scope struct {
	Parent          *Scope
	LexicalLocation Location
	DynamicLocation jsonpointer.T
	Schema          JSON
}

func NewRootScope(location Location, schema JSON) *Scope {
	return &Scope{
		LexicalLocation: location,
		Schema:          schema,
	}
}

func (s *Scope) Descend(t string, schema JSON) *Scope {
	return &Scope{
		Parent:          s,
		LexicalLocation: s.LexicalLocation.AddReferenceToken(t),
		DynamicLocation: s.DynamicLocation.AddReferenceToken(t),
		Schema:          schema,
	}
}

func (s *Scope) Ref(location Location, schema JSON) *Scope {
	return &Scope{
		Parent:          s,
		LexicalLocation: location,
		DynamicLocation: s.DynamicLocation,
		Schema:          schema,
	}
}
