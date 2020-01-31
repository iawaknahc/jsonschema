package jsonschema

import (
	"encoding/json"
	"net/url"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Location struct {
	BaseURI     url.URL
	JSONPointer jsonpointer.T
}

func (l Location) String() string {
	u := l.BaseURI
	u.Fragment = l.JSONPointer.String()
	return u.String()
}

func (l Location) AddReferenceToken(s string) Location {
	return Location{
		BaseURI:     l.BaseURI,
		JSONPointer: l.JSONPointer.AddReferenceToken(s),
	}
}

func (l Location) MarshalJSON() ([]byte, error) {
	u := l.BaseURI
	u.Fragment = l.JSONPointer.String()
	return json.Marshal(u.String())
}

func (l *Location) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	ptr, err := jsonpointer.Parse(u.Fragment)
	if err != nil {
		return err
	}

	u.Fragment = ""
	l.BaseURI = *u
	l.JSONPointer = ptr
	return nil
}
