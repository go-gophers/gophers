package gophers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateRequest(t *testing.T) {
	u, err := url.Parse("https://host.example/prefix/?foo=bar")
	require.Nil(t, err)
	client := NewClient(*u)
	client.DefaultHeaders.Set("X-Header", "123")

	req := client.NewRequest(t, "POST", "/user", nil)
	require.Equal(t, "https://host.example/prefix/user?foo=bar", req.URL.String())
	require.Empty(t, req.RequestURI)
}

func TestColorLoggerFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	now := time.Now().Format(time.RFC3339)
	v := url.Values{}
	v.Add("time", now)

	u, err := url.Parse(server.URL)
	require.Nil(t, err)
	u.RawQuery = v.Encode()

	ft := new(FakeT)
	client := NewClient(*u)
	req := client.NewRequest(t, "GET", "/user", nil)
	client.Do(ft, req, 200)

	require.Equal(t, []string{
		"\n[\x1b[34mGET /user?time=" + url.QueryEscape(now) + " HTTP/1.1\x1b[0m]\n",
		"\n[\x1b[32mHTTP/1.1 200 OK\x1b[0m]\n",
	}, ft.Logs)

	require.Empty(t, ft.Errors)
	require.Empty(t, ft.Fatals)
}

type FakeT struct {
	Logs   []string
	Errors []string
	Fatals []string
}

func (f *FakeT) Logf(format string, a ...interface{}) {
	f.Logs = append(f.Logs, fmt.Sprintf(format, a))
}

func (f *FakeT) Errorf(format string, a ...interface{}) {
	f.Errors = append(f.Errors, fmt.Sprintf(format, a))
}

func (f *FakeT) Fatalf(format string, a ...interface{}) {
	f.Fatals = append(f.Fatals, fmt.Sprintf(format, a))
}
