package jsonschema

type IfThenElse struct{}

var _ Applicator = AllOf{}

func (_ IfThenElse) Apply(ctx ApplicationContext) (annotations []Annotation, errors []Error) {
	if_, ifOK := ctx.GetKeyword("if")
	if !ifOK {
		return
	}

	ifA, ifE := if_.Apply()
	annotations = append(annotations, ifA...)

	if then_, ok := ctx.GetKeyword("then"); ok && len(ifE) <= 0 {
		thenA, thenE := then_.Apply()
		if len(thenE) > 0 {
			annotations = nil
			errors = thenE
		} else {
			annotations = append(annotations, thenA...)
		}
	}
	if else_, ok := ctx.GetKeyword("else"); ok && len(ifE) > 0 {
		elseA, elseE := else_.Apply()
		if len(elseE) > 0 {
			annotations = nil
			errors = elseE
		} else {
			annotations = append(annotations, elseA...)
		}
	}

	return
}
