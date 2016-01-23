package gophers

import (
	"net/url"
	"testing"

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
