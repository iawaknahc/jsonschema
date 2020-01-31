package jsonschema

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

// JSON is a wrapper around a JSON value.
type JSON struct {
	JSONValue         interface{}
	BaseURI           url.URL
	CanonicalLocation jsonpointer.T
}

// WrapJSON wraps plain.
func WrapJSON(plain interface{}) JSON {
	switch v := plain.(type) {
	case map[string]interface{}:
		obj := make(map[string]JSON)
		for key, value := range v {
			obj[key] = WrapJSON(value)
		}
		return JSON{JSONValue: obj}
	case []interface{}:
		arr := make([]JSON, len(v))
		for i, value := range v {
			arr[i] = WrapJSON(value)
		}
		return JSON{JSONValue: arr}
	default:
		return JSON{JSONValue: plain}
	}
}

// UnwrapJSON unwraps j.
func UnwrapJSON(j JSON) interface{} {
	switch v := j.JSONValue.(type) {
	case map[string]JSON:
		obj := make(map[string]interface{})
		for key, value := range v {
			obj[key] = UnwrapJSON(value)
		}
		return obj
	case []JSON:
		arr := make([]interface{}, len(v))
		for i, value := range v {
			arr[i] = UnwrapJSON(value)
		}
		return arr
	default:
		return j.JSONValue
	}
}

// ToFloat64 converts json.Number to float64.
func ToFloat64(plain interface{}) interface{} {
	switch v := plain.(type) {
	case map[string]interface{}:
		obj := make(map[string]interface{})
		for key, value := range v {
			obj[key] = ToFloat64(value)
		}
		return obj
	case []interface{}:
		arr := make([]interface{}, len(v))
		for i, value := range v {
			arr[i] = ToFloat64(value)
		}
		return arr
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			panic(err)
		}
		return f
	default:
		return plain
	}
}

// DecodePlainJSON decodes r into plain JSON with json.Number.
func DecodePlainJSON(r io.Reader) (out interface{}, err error) {
	d := json.NewDecoder(r)
	d.UseNumber()
	err = d.Decode(&out)
	return
}

// TraverseJSON is like jsonpointer.T.Traverse but for JSON.
func TraverseJSON(ptr jsonpointer.T, j JSON) (out *JSON, err error) {
	out = &j
	var i interface{}
	for ptr.More() {
		ptr, i, err = ptr.TraverseOnce(out.JSONValue)
		if err != nil {
			return nil, err
		}
		casted := i.(JSON)
		out = &casted
	}
	return
}
