package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type OutputNode struct {
	// Standard fields
	Valid                   bool          `json:"valid"`
	InstanceLocation        jsonpointer.T `json:"instanceLocation"`
	KeywordLocation         jsonpointer.T `json:"keywordLocation"`
	AbsoluteKeywordLocation Location      `json:"absoluteKeywordLocation"`
	Annotation              interface{}   `json:"annotation,omitempty"`
	Errors                  []OutputNode  `json:"errors,omitempty"`
	Annotations             []OutputNode  `json:"annotations,omitempty"`
	// Extra fields
	Keyword string      `json:"keyword,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}
