package jsonschema

type Keyword interface {
	Keyword() string
	Apply(ctx ApplicationContext, input Node) (*Node, error)
}

type AnnotatingKeyword interface {
	Keyword
	CombineAnnotations(values []interface{}) (interface{}, bool)
}

type Vocabulary struct {
	Keywords []Keyword
}

var DefaultVocabulary = Vocabulary{
	Keywords: []Keyword{
		// Independent keywords
		// Their order is unimportant.
		Type{},
		Const{},
		Enum{},

		MaxItems{},
		MinItems{},
		UniqueItems{},

		MaxLength{},
		MinLength{},
		Pattern{},

		MultipleOf{},
		Maximum{},
		ExclusiveMaximum{},
		Minimum{},
		ExclusiveMinimum{},

		Required{},
		PropertyNames{},
		MaxProperties{},
		MinProperties{},
		// Keywords that must be processed before in-place applicators.
		// The order within the group is also important.
		// The properties group.
		Properties{},
		PatternProperties{},
		AdditionalProperties{},
		// The items group.
		Items{},
		AdditionalItems{},
		// In-place applicators
		AllOf{},
		OneOf{},
		AnyOf{},
		Not{},
		If{},
		Then{},
		Else{},
	},
}
