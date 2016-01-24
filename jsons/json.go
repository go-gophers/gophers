// Package jsons allows sloppy work with JSON structures (objects and arrays).
package jsons

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// TODO support JSONPath - JSON Pointer is not that good

// Struct is common interface for JSON structure.
type Struct interface {
	fmt.Stringer
	Indent() string
	Reader() *strings.Reader
	Get(path string) Struct
	Clone() Struct
	KeepFields(fields ...string) Struct
	RemoveFields(fields ...string) Struct
}

// Object is JSON object structure.
type Object map[string]interface{}

func (j Object) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j Object) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j Object) Reader() *strings.Reader {
	return strings.NewReader(j.String())
}

func (j Object) Get(path string) Struct {
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "" || parts[1] == "" {
		panic(fmt.Errorf("invalid path %q", path))
	}
	f := parts[1]
	s, ok := j[f]
	if !ok {
		panic(fmt.Errorf("key %q not present in object %s (path %q)", f, j.String(), path))
	}
	v := Cast(s)
	if len(parts) == 2 {
		return v
	}
	next := "/" + strings.Join(parts[2:], "/")
	return v.Get(next)
}

func (j Object) Clone() Struct {
	var n Object
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

func (j Object) KeepFields(fields ...string) Struct {
	n := j.Clone().(Object)
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

func (j Object) RemoveFields(fields ...string) Struct {
	n := j.Clone().(Object)
	for _, f := range fields {
		delete(n, f)
	}
	return n
}

// Array is JSON array structure.
type Array []interface{}

func (j Array) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j Array) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (j Array) Reader() *strings.Reader {
	return strings.NewReader(j.String())
}

func (j Array) Clone() Struct {
	var n Array
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

func (j Array) RemoveFields(fields ...string) Struct {
	n := j.Clone().(Array)
	for _, e := range n {
		var o Object
		switch e := e.(type) {
		case Object:
			o = e
		case map[string]interface{}:
			o = Object(e)
		default:
			panic(fmt.Errorf("%v (%T) is not Object", e, e))
		}

		for _, f := range fields {
			delete(o, f)
		}
	}
	return n
}

func (j Array) KeepFields(fields ...string) Struct {
	n := j.Clone().(Array)
	for _, e := range n {
		var o Object
		switch e := e.(type) {
		case Object:
			o = e
		case map[string]interface{}:
			o = Object(e)
		default:
			panic(fmt.Errorf("%v (%T) is not Object", e, e))
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

func (j Array) Get(path string) Struct {
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
	v := Cast(j[n])
	if len(parts) == 2 {
		return v
	}
	next := "/" + strings.Join(parts[2:], "/")
	return v.Get(next)
}

func Cast(v interface{}) Struct {
	switch v := v.(type) {
	case Object:
		return v
	case map[string]interface{}:
		return Object(v)
	case Array:
		return v
	case []interface{}:
		return Array(v)
	default:
		panic(fmt.Errorf("invalid invocation: AsJSON(%v) (%T)", v, v))
	}
}

func Parse(s string, args ...interface{}) Struct {
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
		var o Object
		err := d.Decode(&o)
		if err != nil {
			panic(err)
		}
		return o

	case '[':
		var a Array
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

// check interfaces
var (
	_ Struct = Object{}
	_ Struct = Array{}
)
