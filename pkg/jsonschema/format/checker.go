package format

type FormatChecker interface {
	CheckFormat(value interface{}) error
}

var DefaultChecker map[string]FormatChecker = map[string]FormatChecker{
	"ipv4":         IPV4{},
	"email":        Email{},
	"json-pointer": JSONPointer{},
}
