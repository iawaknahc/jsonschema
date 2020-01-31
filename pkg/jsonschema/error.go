package jsonschema

import (
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Error struct {
	InstanceLocation        jsonpointer.T `json:"instanceLocation"`
	Keyword                 string        `json:"keyword"`
	KeywordLocation         Location      `json:"keywordLocation"`
	AbsoluteKeywordLocation Location      `json:"absoluteKeywordLocation"`
	Value                   interface{}   `json:"value"`
}
