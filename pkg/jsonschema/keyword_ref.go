package jsonschema

type Ref struct{}

var _ Keyword = Ref{}

func (_ Ref) Keyword() string {
	return "$ref"
}

func (_ Ref) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	ref := input.Schema.JSONValue.(string)
	u, err := input.Parent.Schema.BaseURI.Parse(ref)
	if err != nil {
		return nil, err
	}

	referencedSchema, err := ctx.Collection.GetSchema(u.String())
	if err != nil {
		return nil, err
	}
	// TODO(ref): Detect cycle
	// for _, l := range c.ReferencedLocation {
	// 	if l.String() == location.String() {
	// 		return nil, ErrCircularReference{
	// 			Locations: c.ReferencedLocation,
	// 		}
	// 	}
	// }
	// c.ReferencedLocation = append(c.ReferencedLocation, location)
	location := Location{
		BaseURI:     referencedSchema.BaseURI,
		JSONPointer: referencedSchema.CanonicalLocation,
	}
	childInput := Node{
		Valid:                   true,
		Parent:                  input.Parent,
		Instance:                input.Instance,
		InstanceLocation:        input.InstanceLocation,
		Schema:                  *referencedSchema,
		KeywordLocation:         input.KeywordLocation,
		AbsoluteKeywordLocation: location,
	}
	child, err := ctx.Apply(childInput)
	if err != nil {
		return nil, err
	}
	input.Valid = child.Valid
	input.Children = append(input.Children, *child)

	return &input, nil
}
