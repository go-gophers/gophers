package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type JSONValue interface {
	String() string
	Indent() string
	KeepFields(fields ...string) JSONValue
	RemoveFields(fields ...string) JSONValue
}

type JSONObject map[string]interface{}

func (j JSONObject) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j JSONObject) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j JSONObject) KeepFields(fields ...string) JSONValue {
	for k := range j {
		var keep bool
		for _, f := range fields {
			if f == k {
				keep = true
				break
			}
		}
		if !keep {
			delete(j, k)
		}
	}
	return j
}

func (j JSONObject) RemoveFields(fields ...string) JSONValue {
	for _, f := range fields {
		delete(j, f)
	}
	return j
}

type JSONArray []interface{}

func (j JSONArray) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j JSONArray) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j JSONArray) RemoveFields(fields ...string) JSONValue {
	for _, e := range j {
		var o JSONObject
		switch e := e.(type) {
		case JSONObject:
			o = e
		case map[string]interface{}:
			o = JSONObject(e)
		default:
			panic(fmt.Errorf("%v (%T) is not JSONObject", e, e))
		}

		for _, f := range fields {
			delete(o, f)
		}
	}
	return j
}

func (j JSONArray) KeepFields(fields ...string) JSONValue {
	for _, e := range j {
		var o JSONObject
		switch e := e.(type) {
		case JSONObject:
			o = e
		case map[string]interface{}:
			o = JSONObject(e)
		default:
			panic(fmt.Errorf("%v (%T) is not JSONObject", e, e))
		}

		for k := range o {
			var keep bool
			for _, f := range fields {
				if f == k {
					keep = true
					break
				}
			}
			if !keep {
				delete(o, k)
			}
		}
	}
	return j
}

func AsJSON(v interface{}) JSONValue {
	switch v := v.(type) {
	case JSONObject:
		return v
	case map[string]interface{}:
		return JSONObject(v)
	case JSONArray:
		return v
	case []interface{}:
		return JSONArray(v)
	default:
		panic(fmt.Errorf("invalid invocation: AsJSON(%v) (%T)", v, v))
	}
}

func JSON(s string, args ...interface{}) JSONValue {
	if s == "" {
		panic(fmt.Errorf("invalid invocation: JSON(%q)", s))
	}

	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	d := json.NewDecoder(strings.NewReader(s))
	// d.UseNumber()

	switch s[0] {
	case '{':
		var o JSONObject
		err := d.Decode(&o)
		if err != nil {
			panic(err)
		}
		return o

	case '[':
		var a JSONArray
		err := d.Decode(&a)
		if err != nil {
			panic(err)
		}
		return a

	default:
		// TODO handle scalar JSON values?
		panic(fmt.Errorf("unexpected argument: %q", s))
	}
}

func ReadJSON(t testing.TB, r io.Reader) (j JSONValue) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
			j = nil
		}
	}()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
		return
	}

	return JSON(string(b))
}

// check interfaces
var (
	_ JSONValue = JSONObject{}
	_ JSONValue = JSONArray{}
)
