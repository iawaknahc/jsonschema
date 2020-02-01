package jsonschema

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

type Location struct {
	BaseURI     url.URL
	JSONPointer jsonpointer.T
}

func (l Location) String() string {
	l.BaseURI.Fragment = ""
	// Write empty fragment
	return fmt.Sprintf("%s%s", l.BaseURI.String(), l.JSONPointer.Fragment())
}

func (l Location) AddReferenceToken(s string) Location {
	return Location{
		BaseURI:     l.BaseURI,
		JSONPointer: l.JSONPointer.AddReferenceToken(s),
	}
}

func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
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
