package jsonschema

import (
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

// ErrSchemaNotFound occurs when the schema cannot be found.
var ErrSchemaNotFound = errors.New("schema not found")

// ErrFragmentInId occurs when $id contains fragment.
var ErrFragmentInId = errors.New("no fragment is allowed in $id")

// ErrNotASchema occurs when the retrieved schema is not an object or boolean.
var ErrNotASchema = errors.New("not a schema")

// Collection is a collection of schemas that can reference each other.
type Collection struct {
	Index map[string]JSON
}

// NewCollection creates a new Collection.
func NewCollection() *Collection {
	return &Collection{
		Index: map[string]JSON{},
	}
}

// AddSchema adds the schema in r.
// base specifies the root schema $id if $id is missing.
func (c *Collection) AddSchema(r io.Reader, base string) (err error) {
	plainSchema, err := DecodePlainJSON(r)
	schema := WrapJSON(plainSchema)

	baseURL, err := url.Parse(base)
	if err != nil {
		return
	}

	if baseURL.Fragment != "" {
		err = ErrFragmentInId
		return
	}

	_, err = c.buildIndex(schema, *baseURL, nil)
	if err != nil {
		return
	}

	return
}

// GetSchema retrieves the schema referenced by u.
func (c *Collection) GetSchema(u string) (schemaPtr *JSON, err error) {
	// Parse u once to remove empty fragment.
	url, err := url.Parse(u)
	if err != nil {
		return
	}

	var ptr jsonpointer.T

	// If there is a fragment and it starts with a slash,
	// then we find out the schema and follow the JSON Pointer.
	if strings.HasPrefix(url.Fragment, "/") {
		fragment := url.Fragment
		url.Fragment = ""
		ptr, err = jsonpointer.Parse(fragment)
		if err != nil {
			return
		}

		schema, ok := c.Index[url.String()]
		if !ok {
			err = ErrSchemaNotFound
			return
		}

		schemaPtr, err = TraverseJSON(ptr, schema)
		if err != nil {
			return
		}
	} else {
		schema, ok := c.Index[url.String()]
		if !ok {
			err = ErrSchemaNotFound
			return
		}
		schemaPtr = &schema
	}

	switch (*schemaPtr).JSONValue.(type) {
	case map[string]JSON:
		break
	case bool:
		break
	default:
		err = ErrNotASchema
		return
	}

	return
}

func (c *Collection) buildIndex(schema JSON, currentBaseURL url.URL, ptr jsonpointer.T) (result JSON, err error) {
	switch s := schema.JSONValue.(type) {
	case map[string]JSON:
		if id, ok := s["$id"].JSONValue.(string); ok {
			var u *url.URL
			u, err = currentBaseURL.Parse(id)
			if err != nil {
				return
			}
			currentBaseURL = *u

			if currentBaseURL.Fragment != "" {
				err = ErrFragmentInId
				return
			}

			ptr = nil
		}

		schema.BaseURI = currentBaseURL
		schema.CanonicalLocation = ptr

		for key, value := range s {
			value, err = c.buildIndex(value, currentBaseURL, ptr.AddReferenceToken(key))
			if err != nil {
				return
			}
			s[key] = value
		}

		if anchor, ok := s["$anchor"].JSONValue.(string); ok {
			u := currentBaseURL
			u.Fragment = anchor
			c.Index[u.String()] = schema
		}

		result = schema
	case []JSON:
		schema.BaseURI = currentBaseURL
		schema.CanonicalLocation = ptr

		for idx, value := range s {
			value, err = c.buildIndex(value, currentBaseURL, ptr.AddReferenceToken(strconv.Itoa(idx)))
			if err != nil {
				return
			}
			s[idx] = value
		}

		result = schema
	default:
		schema.BaseURI = currentBaseURL
		schema.CanonicalLocation = ptr
		result = schema
	}

	u := currentBaseURL
	u.Fragment = ptr.String()
	c.Index[u.String()] = result

	return
}

// Apply applies the schema referenced by u on r.
func (c *Collection) Apply(u string, r io.Reader) (annotations []Annotation, errors []Error, err error) {
	instance, err := DecodePlainJSON(r)
	if err != nil {
		return
	}

	schema, err := c.GetSchema(u)
	if err != nil {
		return
	}

	location := Location{
		BaseURI:     schema.BaseURI,
		JSONPointer: schema.CanonicalLocation,
	}

	ctx := &ApplicationContext{
		Collection:              c,
		Schema:                  *schema,
		Instance:                instance,
		KeywordLocation:         location,
		AbsoluteKeywordLocation: location,
		Vocabulary:              DefaultVocabulary,
	}

	annotations, errors = ctx.Apply()
	return
}
