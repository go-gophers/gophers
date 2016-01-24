package jsons

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Array is JSON array structure. It implements Struct interface.
type Array []interface{}

// String returns compact JSON representation of JSON array.
// It panics in case of error.
func (j Array) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Indent returns indented JSON representation of JSON array.
// It panics in case of error.
func (j Array) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Get returns JSON substructure by given JSON Pointer path
// (https://tools.ietf.org/html/rfc6901). Scalar values are not supported.
// It panics in case of error.
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

// Clone returns a deep copy of JSON array.
// It panics in case of error.
func (j Array) Clone() Struct {
	var n Array
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

// KeepFields returns a deep copy of JSON array with given object fields kept in all subobjects,
// and all other fields removed.
// It panics in case of error.
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

// RemoveFields returns a deep copy of JSON array with given object fields removed in all subobjects.
// It panics in case of error.
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

// check interfaces
var (
	_ Struct       = Array{}
	_ fmt.Stringer = Array{}
)
