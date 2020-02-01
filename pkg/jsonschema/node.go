package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Node struct {
	// Standard fields
	Valid                   bool          `json:"valid"`
	InstanceLocation        jsonpointer.T `json:"instanceLocation"`
	KeywordLocation         Location      `json:"keywordLocation"`
	AbsoluteKeywordLocation Location      `json:"absoluteKeywordLocation"`
	Annotation              interface{}   `json:"annotation,omitempty"`
	Children                []Node        `json:"children,omitempty"`
	// Extra fields
	Keyword string      `json:"keyword,omitempty"`
	Info    interface{} `json:"info,omitempty"`
	// Runtime fields
	Parent   *Node       `json:"-"`
	Schema   JSON        `json:"-"`
	Instance interface{} `json:"-"`
}

func (n *Node) GetAnnotationsFromAdjacentKeywords(k AnnotatingKeyword) (interface{}, bool) {
	if n.Parent == nil {
		return nil, false
	}

	var values []interface{}
	keyword := k.Keyword()
	for _, child := range n.Parent.Children {
		if n.InstanceLocation.String() == child.InstanceLocation.String() && child.Keyword == keyword {
			values = append(values, child.Annotation)
		}
	}

	return k.CombineAnnotations(values)
}
