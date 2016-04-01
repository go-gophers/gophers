package gophers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

func TestColorLoggerFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	now := time.Now().Format(time.RFC3339)
	v := url.Values{}
	v.Add("time", now)

	u, err := url.Parse(server.URL)
	require.Nil(t, err)
	u.RawQuery = v.Encode()

	tb := new(FakeTB)
	client := NewClient(*u)
	req := client.NewRequest(tb, "GET", "/user", nil)
	client.Do(tb, req, 200)

	assert.Equal(tb, []string{
		"\n[\x1b[34mGET /user?time=" + url.QueryEscape(now) + " HTTP/1.1\x1b[0m]\n",
		"\n[\x1b[32mHTTP/1.1 200 OK\x1b[0m]\n",
	}, tb.Logs)
	assert.Empty(tb, tb.Errors)
	assert.Empty(tb, tb.Fatals)
}
