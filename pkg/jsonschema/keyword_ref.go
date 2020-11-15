package jsonschema

type Ref struct{}

var _ Keyword = Ref{}

func (_ Ref) Keyword() string {
	return "$ref"
}

func (_ Ref) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	ref := input.Scope.Schema.JSONValue.(string)
	u, err := input.Parent.Scope.Schema.BaseURI.Parse(ref)
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
	childInput := Node{
		Valid:            true,
		Parent:           input.Parent,
		Instance:         input.Instance,
		InstanceLocation: input.InstanceLocation,
		Scope: input.Scope.Ref(Location{
			BaseURI:     referencedSchema.BaseURI,
			JSONPointer: referencedSchema.CanonicalLocation,
		}, *referencedSchema),
	}
	child, err := ctx.Apply(childInput)
	if err != nil {
		return nil, err
	}
	input.Valid = child.Valid
	input.Children = append(input.Children, *child)

	return &input, nil
}
