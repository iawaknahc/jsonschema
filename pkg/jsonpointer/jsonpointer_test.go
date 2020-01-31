package jsonpointer

import (
	"encoding/json"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input  string
		errMsg string
	}{
		{"", ""},
		{"/", ""},
		{"//", ""},
		{"/a", ""},
		{"/a/", ""},
		{"/a//", ""},
		{"/a/b", ""},
		{"/a/b/", ""},
		{"/a/b//", ""},
		{"/a/b//a/b", ""},
		{"/~0", ""},
		{"/~0~1", ""},
		{"/~0~1/~1~0", ""},
		{"/~0~1/~1~0/~0~1", ""},

		{"a", `0: expecting / but found: "a"`},
		{"/~", `1: expecting 0 or 1 but found: EOF`},
		{"/~2", `2: expecting 0 or 1 but found: "2"`},
		{"/a/b/~3", `6: expecting 0 or 1 but found: "3"`},
	}
	for _, c := range cases {
		actual, err := Parse(c.input)
		if err != nil {
			if c.errMsg != err.Error() {
				t.Errorf("%s != %s", err.Error(), c.errMsg)
			}
		} else if c.errMsg != "" {
			t.Errorf("expecting error: %s", c.errMsg)
		} else {
			str := actual.String()
			if str != c.input {
				t.Errorf("%s != %s", str, c.input)
			}
		}
	}
}

func TestParseFragment(t *testing.T) {
	cases := []struct {
		input  string
		errMsg string
	}{
		{"#", ""},
		{"#/", ""},
		{"#//", ""},
		{"#/a", ""},
		{"#/a/", ""},
		{"#/a//", ""},
		{"#/a/b", ""},
		{"#/a/b/", ""},
		{"#/a/b//", ""},
		{"#/a/b//a/b", ""},
		{"#/~0", ""},
		{"#/~0~1", ""},
		{"#/~0~1/~1~0", ""},
		{"#/~0~1/~1~0/~0~1", ""},
	}
	for _, c := range cases {
		actual, err := Parse(c.input)
		if err != nil {
			if c.errMsg != err.Error() {
				t.Errorf("%s != %s", err.Error(), c.errMsg)
			}
		} else if c.errMsg != "" {
			t.Errorf("expecting error: %s", c.errMsg)
		} else {
			str := actual.Fragment()
			if str != c.input {
				t.Errorf("%s != %s", str, c.input)
			}
		}
	}
}

func TestAddReferenceToken(t *testing.T) {
	cases := []struct {
		input    string
		refToken string
		expected string
	}{
		{"", "a", "/a"},
		{"/a", "a", "/a/a"},
		{"", "", "/"},
		{"/", "", "//"},
		{"/a", "", "/a/"},
		{"", "~", "/~0"},
		{"", "/", "/~1"},
		{"", "~//~", "/~0~1~1~0"},
	}
	for _, c := range cases {
		ptr, err := Parse(c.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else {
			newPtr := ptr.AddReferenceToken(c.refToken)
			if newPtr.String() != c.expected {
				t.Errorf("%s != %s", ptr.String(), c.expected)
			}
			if ptr.String() != c.input {
				t.Errorf("AddReferenceToken should be immutable")
			}
		}
	}
}

func TestTraverse(t *testing.T) {
	cases := []struct {
		value    string
		ptr      string
		expected string
	}{
		{
			`null`,
			"",
			`null`,
		},
		{
			`{"a": 1}`,
			"/a",
			`1`,
		},
		{
			`["a", "b"]`,
			"/1",
			`"b"`,
		},
		{
			`{"a": [{"a": 1}]}`,
			"/a/0/a",
			`1`,
		},
		{
			`{"a": [{"a": 1}]}`,
			"/a/0",
			`{"a":1}`,
		},
	}
	for _, c := range cases {
		ptr, err := Parse(c.ptr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else {
			var value interface{}
			err := json.Unmarshal([]byte(c.value), &value)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else {
				traversed, err := ptr.Traverse(value)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					marshaled, _ := json.Marshal(traversed)
					if string(marshaled) != c.expected {
						t.Errorf("%s != %s", string(marshaled), c.expected)
					}
				}
			}
		}
	}
}

func TestFragment(t *testing.T) {
	cases := []struct {
		str  string
		frag string
	}{
		{"", "#"},
		{"/foo", "#/foo"},
		{"/foo/0", "#/foo/0"},
		{"/", "#/"},
		{"/a~1b", "#/a~1b"},
		{"/c%d", "#/c%25d"},
		{"/e^f", "#/e%5Ef"},
		{"/g|h", "#/g%7Ch"},
		{"/h\\j", "#/h%5Cj"},
		{"/k\"l", "#/k%22l"},
		{"/ ", "#/%20"},
		{"/m~0n", "#/m~0n"},
	}
	for _, c := range cases {
		escaped := FragmentEscape(c.str)
		if escaped != c.frag {
			t.Errorf("%s != %s", escaped, c.frag)
		}

		unescaped, err := FragmentUnescape(c.frag)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if unescaped != c.str {
			t.Errorf("%s != %s", unescaped, c.str)
		}
	}
}

func TestZeroPointer(t *testing.T) {
	var ptr T

	if ptr.String() != "" {
		t.Errorf("zero String() should be empty string")
	}

	if ptr.Fragment() != "#" {
		t.Errorf("zero Fragment() should be #")
	}

	if ptr.AddReferenceToken("foo").String() != "/foo" {
		t.Errorf("zero AddReferenceToken should be ok")
	}

	if ptr.More() != false {
		t.Errorf("zero More() should be false")
	}

	input := 1
	output, err := ptr.Traverse(input)
	if input != output {
		t.Errorf("zero Traverse should do nothing")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMarshalJSON(t *testing.T) {
	type A struct {
		Ptr T `json:"ptr"`
	}

	a := A{
		Ptr: MustParse("/a/b/c"),
	}

	b, err := json.Marshal(a)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var aa A
	err = json.Unmarshal(b, &aa)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if a.Ptr.String() != aa.Ptr.String() {
		t.Errorf("%s != %s", a.Ptr.String(), aa.Ptr.String())
	}
}
