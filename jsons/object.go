package jsons

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Object is JSON object structure. It implements Struct interface.
type Object map[string]interface{}

// String returns compact JSON representation of JSON object.
// It panics in case of error.
func (j Object) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Indent returns indented JSON representation of JSON object.
// It panics in case of error.
func (j Object) Indent() string {
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Get returns JSON substructure by given JSON Pointer path
// (https://tools.ietf.org/html/rfc6901). Scalar values are not supported.
// It panics in case of error.
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

// Clone returns a deep copy of JSON object.
// It panics in case of error.
func (j Object) Clone() Struct {
	var n Object
	err := json.Unmarshal([]byte(j.String()), &n)
	if err != nil {
		panic(err)
	}
	return n
}

// KeepFields returns a deep copy of JSON object with given fields kept,
// and all other fields removed.
// It panics in case of error.
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

// RemoveFields returns a deep copy of JSON object with given fields removed.
// It panics in case of error.
func (j Object) RemoveFields(fields ...string) Struct {
	n := j.Clone().(Object)
	for _, f := range fields {
		delete(n, f)
	}
	return n
}

// check interfaces
var (
	_ Struct       = Object{}
	_ fmt.Stringer = Object{}
)
