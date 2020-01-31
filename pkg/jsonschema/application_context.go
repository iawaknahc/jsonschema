package jsonschema

import (
	"reflect"
	"strings"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
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
	Collection              *Collection
	Schema                  JSON
	Instance                interface{}
	Keyword                 string
	InstanceLocation        jsonpointer.T
	KeywordLocation         Location
	AbsoluteKeywordLocation Location
	Vocabulary              Vocabulary
	Annotations             []Annotation
	ReferencedLocation      []Location
}

func (c ApplicationContext) Apply() (annotations []Annotation, errors []Error) {
	c.Annotations = nil

	if schema, ok := c.Schema.JSONValue.(map[string]JSON); ok {
		if ref, ok := schema["$ref"].JSONValue.(string); ok {
			u, err := c.Schema.BaseURI.Parse(ref)
			if err != nil {
				panic(err)
			}
			referencedSchema, err := c.Collection.GetSchema(u.String())
			if err != nil {
				panic(err)
			}

			c.Schema = *referencedSchema
			c.KeywordLocation = c.KeywordLocation.AddReferenceToken("$ref")
			location := Location{
				BaseURI:     referencedSchema.BaseURI,
				JSONPointer: referencedSchema.CanonicalLocation,
			}

			for _, l := range c.ReferencedLocation {
				if l.String() == location.String() {
					panic(ErrCircularReference{
						Locations: c.ReferencedLocation,
					})
				}
			}
			c.ReferencedLocation = append(c.ReferencedLocation, location)

			c.AbsoluteKeywordLocation = location
			return c.Apply()
		}
	}

	c.ReferencedLocation = nil

	for _, applicator := range c.Vocabulary.Independent {
		a, e := c.ApplyApplicator(applicator)
		c.Annotations = append(c.Annotations, a...)
		errors = append(errors, e...)
	}

	for _, applicator := range c.Vocabulary.PreInplace {
		a, e := c.ApplyApplicator(applicator)
		c.Annotations = append(c.Annotations, a...)
		errors = append(errors, e...)
	}

	for _, applicator := range c.Vocabulary.Inplace {
		a, e := c.ApplyApplicator(applicator)
		c.Annotations = append(c.Annotations, a...)
		errors = append(errors, e...)
	}

	for _, applicator := range c.Vocabulary.PostInplace {
		a, e := c.ApplyApplicator(applicator)
		c.Annotations = append(c.Annotations, a...)
		errors = append(errors, e...)
	}

	annotations = c.Annotations

	return
}

func (c ApplicationContext) ApplyApplicator(applicator Applicator) (annotations []Annotation, errors []Error) {
	keyworder, ok := applicator.(Keyworder)
	if !ok {
		return applicator.Apply(c)
	}

	keyword := keyworder.Keyword()
	ctx, ok := c.GetKeyword(keyword)
	if !ok {
		return
	}

	return applicator.Apply(*ctx)
}

func (c ApplicationContext) GetKeyword(keyword string) (*ApplicationContext, bool) {
	if schema, ok := c.Schema.JSONValue.(map[string]JSON); ok {
		if j, ok := schema[keyword]; ok {
			c.Schema = j
			c.Keyword = keyword
			c.KeywordLocation = c.KeywordLocation.AddReferenceToken(keyword)
			c.AbsoluteKeywordLocation = c.AbsoluteKeywordLocation.AddReferenceToken(keyword)
			return &c, true
		}
	}
	return nil, false
}

func (c ApplicationContext) GetAnnotation(annotator Annotator) (*Annotation, bool) {
	keyword := annotator.Keyword()
	var annotations []Annotation
	for _, a := range c.Annotations {
		if reflect.DeepEqual(c.InstanceLocation, a.InstanceLocation) && keyword == a.Keyword {
			annotations = append(annotations, a)
		}
	}
	if len(annotations) <= 0 {
		return nil, false
	} else if len(annotations) == 1 {
		return &annotations[0], true
	}
	return annotator.MergeAnnotations(annotations)
}
