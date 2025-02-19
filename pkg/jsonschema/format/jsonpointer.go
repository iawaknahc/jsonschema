package format

import (
	"context"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type JSONPointer struct{}

var _ FormatChecker = JSONPointer{}

func (JSONPointer) CheckFormat(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	_, err := jsonpointer.ParseStringRepresentation(str)
	if err != nil {
		return err
	}

	return nil
}
