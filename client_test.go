package gophers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-gophers/gophers/jsons"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRequest(t *testing.T) {
	u, err := url.Parse("https://host.example/prefix/?foo=bar")
	require.Nil(t, err)
	client := NewClient(*u)
	client.DefaultHeaders.Set("X-Header", "123")

	req := client.NewRequest(t, "POST", "/user", nil)
	assert.Equal(t, "https://host.example/prefix/user?foo=bar", req.URL.String())
	assert.Empty(t, req.RequestURI)
	assert.Equal(t, http.Header{"User-Agent": {defaultUserAgent}, "X-Header": {"123"}}, req.Header)
}

func TestRequestResponseBody(t *testing.T) {
	u, err := url.Parse("http://jsonplaceholder.typicode.com")
	require.Nil(t, err)
	client := NewClient(*u)

	j := jsons.Parse(`{"userId": 1, "id": 101, "title": "title", "body": "body"}`)
	req := client.NewRequest(t, "POST", "/posts", j)
	assert.Nil(t, req.Body)
	assert.NotNil(t, req.Request.Body)

	resp := client.Do(t, req, 201)
	assert.Equal(t, []byte(j.String()), req.Body)
	assert.IsType(t, errorReadCloser{}, req.Request.Body)
	assert.Equal(t, jsons.Parse(`{"id": 101}`), jsons.ParseBytes(resp.Body))
	assert.IsType(t, errorReadCloser{}, resp.Response.Body)
}
