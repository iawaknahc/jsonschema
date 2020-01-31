package jsonpointer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// ErrInvalidFragment occurs when the given input is not a URL fragment.
var ErrInvalidFragment = errors.New("invalid fragment")

type state int

const (
	stateSlash state = iota
	stateEscape
	stateChar
)

// T is a JSON pointer.
// A JSON pointer is a sequence of reference tokens.
type T []string

// Parse parses the input into T.
// input can either in string representationn or URL fragment representation.
func Parse(input string) (T, error) {
	var err error
	if strings.HasPrefix(input, "#") {
		input, err = FragmentUnescape(input)
		if err != nil {
			return nil, err
		}
	}

	var output []string
	var w *strings.Builder
	state := stateSlash

	var idx int
	for i, r := range input {
		idx = i
		switch state {
		case stateSlash:
			if r != '/' {
				return nil, fmt.Errorf("%d: expecting / but found: %#v", idx, string(r))
			}
			w = &strings.Builder{}
			state = stateChar
		case stateEscape:
			switch r {
			case '0':
				w.WriteRune('~')
			case '1':
				w.WriteRune('/')
			default:
				return nil, fmt.Errorf("%d: expecting 0 or 1 but found: %#v", idx, string(r))
			}
			state = stateChar
		case stateChar:
			switch r {
			case '~':
				state = stateEscape
			case '/':
				output = append(output, w.String())
				w.Reset()
			default:
				w.WriteRune(r)
			}
		}
	}

	if state == stateEscape {
		return nil, fmt.Errorf("%d: expecting 0 or 1 but found: EOF", idx)
	}

	if w != nil {
		output = append(output, w.String())
	}

	return T(output), nil
}

// MustParse is a shorthand for creating literal value of T.
func MustParse(input string) T {
	t, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return t
}

// String returns the string representation.
func (t T) String() string {
	w := &strings.Builder{}
	for _, referenceToken := range t {
		w.WriteRune('/')
		for _, r := range referenceToken {
			switch r {
			case '~':
				w.WriteString("~0")
			case '/':
				w.WriteString("~1")
			default:
				w.WriteRune(r)
			}
		}
	}
	return w.String()
}

// MarshalJSON implements Marshaler.
func (t T) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Fragment())
}

// UnmarshalJSON implements Unmarshaler.
func (t *T) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*t = parsed
	return nil
}

// Fragment returns the URL fragment representation.
func (t T) Fragment() string {
	return FragmentEscape(t.String())
}

// AddReferenceToken returns a new JSON pointer with s added.
func (t T) AddReferenceToken(s string) T {
	output := make([]string, len(t))
	copy(output, t)
	output = append(output, s)
	return output
}

// More tells if t has any reference tokens.
func (t T) More() bool {
	return len(t) > 0
}

// TraverseOnce consumes the first reference token.
// v can either be map with string key or any slice.
func (t T) TraverseOnce(v interface{}) (T, interface{}, error) {
	head := t[0]
	tail := t[1:]
	reflectValue := reflect.ValueOf(v)
	switch reflectValue.Kind() {
	case reflect.Map:
		iter := reflectValue.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()
			if key.String() == head {
				return tail, val.Interface(), nil
			}
		}
		return t, nil, fmt.Errorf("%v is undefined", head)
	case reflect.Slice:
		index, err := strconv.Atoi(head)
		if err != nil {
			return t, nil, fmt.Errorf("%v is not an index", head)
		}
		if index < 0 || index >= reflectValue.Len() {
			return t, nil, fmt.Errorf("%v is not in range [0,%d]", index, reflectValue.Len())
		}
		return tail, reflectValue.Index(index).Interface(), nil
	default:
		return t, nil, fmt.Errorf("%v (%T) is not traversable", v, v)
	}
}

// Traverse traverses v with t.
func (t T) Traverse(v interface{}) (out interface{}, err error) {
	out = v
	for t.More() {
		t, out, err = t.TraverseOnce(out)
		if err != nil {
			return nil, err
		}
	}
	return
}

// FragmentEscape is like url.QueryEscape but for fragment.
func FragmentEscape(s string) string {
	// Special case
	if s == "" {
		return "#"
	}
	u := &url.URL{
		Fragment: s,
	}
	return u.String()
}

// FragmentUnescape is like url.QueryUnescape but for fragment.
func FragmentUnescape(s string) (out string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	expected := url.URL{
		Fragment: u.Fragment,
	}
	if *u != expected {
		err = ErrInvalidFragment
		return
	}
	out = u.Fragment
	return
}
