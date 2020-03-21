package jsonschema

type OutputNode struct {
	// Standard fields
	Valid                   bool         `json:"valid"`
	InstanceLocation        string       `json:"instanceLocation"`
	KeywordLocation         string       `json:"keywordLocation"`
	AbsoluteKeywordLocation string       `json:"absoluteKeywordLocation,omitempty"`
	Annotation              interface{}  `json:"annotation,omitempty"`
	Errors                  []OutputNode `json:"errors,omitempty"`
	Annotations             []OutputNode `json:"annotations,omitempty"`
	// Extra fields
	Keyword string      `json:"keyword,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}
