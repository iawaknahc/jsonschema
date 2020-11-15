package jsonschema

import (
	"fmt"
)

type RecursiveRef struct{}

var _ Keyword = RecursiveRef{}

func (_ RecursiveRef) Keyword() string {
	return "$recursiveRef"
}

func (_ RecursiveRef) Apply(ctx ApplicationContext, input Node) (*Node, error) {
	ref := input.Scope.Schema.JSONValue.(string)
	if ref != "#" {
		return nil, fmt.Errorf("$recursiveRef is only defined for the value '#': %v", input.Scope.DynamicLocation)
	}

	// $recursiveRef behave in the same manner as $ref
	u, err := input.Parent.Scope.Schema.BaseURI.Parse(ref)
	if err != nil {
		return nil, err
	}

	referencedSchema, err := ctx.Collection.GetSchema(u.String())
	if err != nil {
		return nil, err
	}

	if obj, ok := referencedSchema.JSONValue.(map[string]JSON); ok {
		if recursiveAnchor, ok := obj["$recursiveAnchor"].JSONValue.(bool); ok && recursiveAnchor {
			outermost := input.Scope

			curr := input.Scope
			for curr != nil {
				if obj, ok := curr.Schema.JSONValue.(map[string]JSON); ok {
					if recursiveAnchor, ok := obj["$recursiveAnchor"].JSONValue.(bool); ok && recursiveAnchor {
						outermost = curr
					}
				}
				curr = curr.Parent
			}

			u, err = outermost.Schema.BaseURI.Parse(ref)
			if err != nil {
				return nil, err
			}
		}

		referencedSchema, err = ctx.Collection.GetSchema(u.String())
		if err != nil {
			return nil, err
		}
	}

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
