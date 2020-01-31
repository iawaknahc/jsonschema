package jsonschema

type Applicator interface {
	Apply(c ApplicationContext) (annotations []Annotation, errors []Error)
}

type Keyworder interface {
	Keyword() string
}

type Vocabulary struct {
	Independent []Applicator
	PreInplace  []Applicator
	Inplace     []Applicator
	PostInplace []Applicator
}

var DefaultVocabulary = Vocabulary{
	Independent: []Applicator{
		Boolean{},

		Type{},

		Const{},

		Required{},

		MaxItems{},
		MinItems{},

		MaxLength{},
		MinLength{},

		MultipleOf{},
		Maximum{},
		ExclusiveMaximum{},
		Minimum{},
		ExclusiveMinimum{},
	},
	PreInplace: []Applicator{
		Properties{},
		PatternProperties{},
		AdditionalProperties{},

		Items{},
		AdditionalItems{},
	},
	Inplace: []Applicator{
		AllOf{},
		OneOf{},
		AnyOf{},
		Not{},
		IfThenElse{},
	},
}
