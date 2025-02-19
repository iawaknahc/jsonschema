package jsonschema

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
	"github.com/iawaknahc/jsonschema/pkg/jsonschema/format"
)

// ErrSchemaNotFound occurs when the schema cannot be found.
var ErrSchemaNotFound = errors.New("schema not found")

// ErrFragmentInId occurs when $id contains fragment.
var ErrFragmentInId = errors.New("no fragment is allowed in $id")

// ErrNotASchema occurs when the retrieved schema is not an object or boolean.
var ErrNotASchema = errors.New("not a schema")

// ErrMetaschemaInSubschema occurs when $schema appears in a subschema.
var ErrMetaschemaInSubschema = errors.New("$schema in not allowed in subschema")

const MetaschemaURI = "https://json-schema.org/draft/2019-09/schema"

// Collection is a collection of schemas that can reference each other.
type Collection struct {
	Index         map[string]JSON
	FormatChecker map[string]format.FormatChecker
}

// NewCollection creates an empty Collection.
func NewCollection() *Collection {
	checker := map[string]format.FormatChecker{}
	for k, v := range format.DefaultChecker {
		checker[k] = v
	}
	return &Collection{
		Index:         map[string]JSON{},
		FormatChecker: checker,
	}
}

// NewMetaschemaCollection creates a Collection that can be used to validate schemas.
func NewMetaschemaCollection() *Collection {
	checker := map[string]format.FormatChecker{}
	for k, v := range format.DefaultChecker {
		checker[k] = v
	}
	c := &Collection{
		Index:         map[string]JSON{},
		FormatChecker: checker,
	}

	var err error

	err = c.AddSchema(strings.NewReader(MetaschemaJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaCoreJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaApplicatorJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaValidationJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaMetadataJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaFormatJSONString), "")
	if err != nil {
		panic(err)
	}

	err = c.AddSchema(strings.NewReader(MetaContentJSONString), "")
	if err != nil {
		panic(err)
	}

	return c
}

// AddSchema adds the schema in r.
// base specifies the root schema $id if $id is missing.
func (c *Collection) AddSchema(r io.Reader, base string) (err error) {
	// FIXME: Support adding metaschema, e.g. validate or ignore $vocabulary
	plainSchema, err := DecodePlainJSON(r)
	if err != nil {
		return
	}
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
		if metaSchema, ok := s["$schema"].JSONValue.(string); ok {
			if len(ptr) > 0 {
				err = ErrMetaschemaInSubschema
				return
			}
			if metaSchema != MetaschemaURI {
				err = fmt.Errorf("unsupported $schema: %v", metaSchema)
				return
			}
		}

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
func (c *Collection) Apply(ctx context.Context, u string, r io.Reader) (node *Node, err error) {
	instance, err := DecodePlainJSON(r)
	if err != nil {
		return
	}

	schema, err := c.GetSchema(u)
	if err != nil {
		return
	}

	appCtx := ApplicationContext{
		Context:      ctx,
		Collection:   c,
		Vocabulary:   DefaultVocabulary,
		PatternCache: &sync.Map{},
	}

	input := Node{
		Valid:    true,
		Instance: instance,
		Scope: NewRootScope(Location{
			BaseURI:     schema.BaseURI,
			JSONPointer: schema.CanonicalLocation,
		}, *schema),
	}

	return appCtx.Apply(input)
}
