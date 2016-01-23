package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

// TODO support JSONPath - JSON Pointer is not that good

type JSONValue interface {
	String() string
	Indent() string
	Get(path string) JSONValue
	Clone() JSONValue
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

func (j JSONObject) Get(path string) JSONValue {
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "" || parts[1] == "" {
		panic(fmt.Errorf("invalid path %q", path))
	}
	f := parts[1]
	s, ok := j[f]
	if !ok {
		panic(fmt.Errorf("key %q not present in object %s (path %q)", f, j.String(), path))
	}
	v := AsJSON(s)
	if len(parts) == 2 {
		return v
	}
	next := "/" + strings.Join(parts[2:], "/")
	return v.Get(next)
}

func (j JSONObject) Clone() JSONValue {
	var n JSONObject
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

func (j JSONObject) KeepFields(fields ...string) JSONValue {
	n := j.Clone().(JSONObject)
	for k := range n {
		var keep bool
		for _, f := range fields {
			if f == k {
				keep = true
				break
			}
		}
		if !keep {
			delete(n, k)
		}
	}
	return n
}

func (j JSONObject) RemoveFields(fields ...string) JSONValue {
	n := j.Clone().(JSONObject)
	for _, f := range fields {
		delete(n, f)
	}
	return n
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

func (j JSONArray) Clone() JSONValue {
	var n JSONArray
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

func (j JSONArray) RemoveFields(fields ...string) JSONValue {
	n := j.Clone().(JSONArray)
	for _, e := range n {
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
	return n
}

func (j JSONArray) KeepFields(fields ...string) JSONValue {
	n := j.Clone().(JSONArray)
	for _, e := range n {
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
	return n
}

func (j JSONArray) Get(path string) JSONValue {
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "" || parts[1] == "" {
		panic(fmt.Errorf("invalid path %q", path))
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Errorf("invalid path %q (%s)", path, err))
	}
	if n >= len(j) {
		panic(fmt.Errorf("index %d not present in array %s (path %q)", n, j.String(), path))
	}
	v := AsJSON(j[n])
	if len(parts) == 2 {
		return v
	}
	next := "/" + strings.Join(parts[2:], "/")
	return v.Get(next)
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
