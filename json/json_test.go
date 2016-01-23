package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONObject(t *testing.T) {
	v := JSON(`{"foo": "bar", "baz": 42, "argument": %t}`, true)
	assert.Equal(t, `{"argument":true,"baz":42,"foo":"bar"}`, v.String())
	indent := `{
  "argument": true,
  "baz": 42,
  "foo": "bar"
}`
	assert.Equal(t, indent, v.Indent())
	assert.Equal(t, JSON(`{"argument": true, "baz": 42}`), v.RemoveFields("nonexisting", "foo"))
	assert.Equal(t, JSON(`{"baz":42}`), v.KeepFields("nonexisting", "baz"))
}

func TestJSONArray(t *testing.T) {
	v := JSON(`[{"foo": "bar1"}, {"foo": "bar2", "baz": 42}, {"foo": "bar3", "argument": %t}]`, true)
	assert.Equal(t, `[{"foo":"bar1"},{"baz":42,"foo":"bar2"},{"argument":true,"foo":"bar3"}]`, v.String())
	indent := `[
  {
    "foo": "bar1"
  },
  {
    "baz": 42,
    "foo": "bar2"
  },
  {
    "argument": true,
    "foo": "bar3"
  }
]`
	assert.Equal(t, indent, v.Indent())
	assert.Equal(t, JSON(`[{},{"baz":42},{"argument":true}]`), v.RemoveFields("nonexisting", "foo"))
	assert.Equal(t, JSON(`[{},{"baz":42},{}]`), v.KeepFields("nonexisting", "baz"))
}

func TestJSONPointer(t *testing.T) {
	v := JSON(`{"foo": [ 0, 1, {"baz": ["good"]} ]}`)
	assert.Equal(t, JSON(`["good"]`), v.Get("/foo/2/baz"))
}

// TODO add tests from RFC, handle escaping
