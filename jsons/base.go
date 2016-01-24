// Package jsons allows sloppy work with JSON structures (objects and arrays).
package jsons

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TODO support JSONPath - JSON Pointer is not that good

// Struct is common interface for JSON structure.
type Struct interface {
	// String returns compact JSON representation of JSON structure.
	String() string

	// Indent returns indented JSON representation of JSON structure.
	Indent() string

	// Get returns JSON substructure by given JSON Pointer path
	// (https://tools.ietf.org/html/rfc6901). Scalar values are not supported.
	Get(path string) Struct

	// Clone returns a deep copy of JSON structure.
	Clone() Struct

	// KeepFields returns a deep copy of JSON structure with given object fields keeps,
	// and all other removed.
	KeepFields(fields ...string) Struct

	// RemoveFields returns a deep copy of JSON structure with given object fields removed.
	RemoveFields(fields ...string) Struct
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
