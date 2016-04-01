package jsons

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Struct is common interface for JSON structures.
type Struct interface {
	// String returns compact JSON representation of JSON structure.
	// It panics in case of error.
	String() string

	// Indent returns indented JSON representation of JSON structure.
	// It panics in case of error.
	Indent() string

	// Get returns JSON substructure by given JSON Pointer path
	// (https://tools.ietf.org/html/rfc6901). Scalar values are not supported.
	// It panics in case of error.
	Get(path string) Struct

	// Clone returns a deep copy of JSON structure.
	// It panics in case of error.
	Clone() Struct

	// KeepFields returns a deep copy of JSON structure with given object fields kept,
	// and all other fields removed.
	// It panics in case of error.
	KeepFields(fields ...string) Struct

	// RemoveFields returns a deep copy of JSON structure with given object fields removed.
	// It panics in case of error.
	RemoveFields(fields ...string) Struct
}

// TODO support JSONPath - JSON Pointer is not that good

// Cast makes type assertion for given value and retuns it as JSON structure.
// It supports Array / []interface{} and Object / map[string]interface{}.
// It panics for other types.
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

// Parse makes JSON structure from given JSON string with fmt verbs and args.
// Scalar values are not supported.
// It panics in case of error.
func Parse(s string, args ...interface{}) Struct {
	if s == "" {
		panic(fmt.Errorf("invalid invocation: JSON(%q)", s))
	}

	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	d := json.NewDecoder(strings.NewReader(s))

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

// ParseBytes is a convenience function to Parse bytes.
func ParseBytes(b []byte) Struct {
	return Parse(string(b))
}
