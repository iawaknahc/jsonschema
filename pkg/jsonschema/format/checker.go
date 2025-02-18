package format

import "context"

type FormatChecker interface {
	CheckFormat(ctx context.Context, value interface{}) error
}

var DefaultChecker map[string]FormatChecker = map[string]FormatChecker{
	"ipv4":         IPV4{},
	"email":        Email{},
	"json-pointer": JSONPointer{},
}
