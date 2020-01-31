package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Annotator interface {
	Keyworder
	MergeAnnotations(annotations []Annotation) (*Annotation, bool)
}

type Annotation struct {
	InstanceLocation        jsonpointer.T `json:"instanceLocation"`
	Keyword                 string        `json:"keyword"`
	KeywordLocation         Location      `json:"keywordLocation"`
	AbsoluteKeywordLocation Location      `json:"absoluteKeywordLocation"`
	Value                   interface{}   `json:"value"`
}
