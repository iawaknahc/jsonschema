package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Node struct {
	// Standard fields
	Valid                   bool
	InstanceLocation        jsonpointer.T
	KeywordLocation         jsonpointer.T
	AbsoluteKeywordLocation Location
	Annotation              interface{}
	Children                []Node
	// Extra fields
	Keyword string
	Info    interface{}
	// Runtime fields
	Parent   *Node
	Schema   JSON
	Instance interface{}
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

func (n *Node) Verbose() (out OutputNode) {
	out = OutputNode{
		Valid:                   n.Valid,
		InstanceLocation:        n.InstanceLocation.Fragment(),
		KeywordLocation:         n.KeywordLocation.Fragment(),
		AbsoluteKeywordLocation: n.AbsoluteKeywordLocation.String(),
		Annotation:              n.Annotation,
		Keyword:                 n.Keyword,
		Info:                    n.Info,
	}
	if n.Valid {
		for _, child := range n.Children {
			out.Annotations = append(out.Annotations, child.Verbose())
		}
	} else {
		for _, child := range n.Children {
			out.Errors = append(out.Errors, child.Verbose())
		}
	}
	return out
}
