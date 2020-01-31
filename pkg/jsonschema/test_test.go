package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"
)

func JSONEq(t *testing.T, expectedStr string, actual interface{}) {
	b, err := json.Marshal(actual)
	if err != nil {
		t.Errorf("JSONEq: unexpected err: %v", err)
	}

	var expected interface{}
	err = json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Errorf("JSONEq: unexpected err: %v", err)
	}

	var actualJ interface{}
	err = json.Unmarshal(b, &actualJ)
	if err != nil {
		t.Errorf("JSONEq: unexpected err: %v", err)
	}

	if !reflect.DeepEqual(actualJ, expected) {
		t.Errorf("%s != %s", string(b), expectedStr)
	}
}
