package github

import (
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/antonholmquist/jason"
	"github.com/stretchr/testify/require"

	"github.com/gophergala2016/gophers"
)

var (
	Client = gophers.NewClient(*gophers.MustParseURL("https://api.github.com/?access_token=" + os.Getenv("GOPHERS_GITHUB_TOKEN")))
)

func TestListOrgs(t *testing.T) {
	req := Client.NewRequest(t, "GET", "/user/orgs")
	b, err := httputil.DumpRequestOut(req, true)
	require.Nil(t, err)
	t.Logf("Request:\n%s", b)

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.Nil(t, err)
	b, err = httputil.DumpResponse(resp, true)
	require.Nil(t, err)
	t.Logf("Response:\n%s", b)

	v, err := jason.NewValueFromReader(resp.Body)
	require.Nil(t, err)
	a, err := v.Array()
	require.Nil(t, err)

	var found bool
	for _, v = range a {
		o, err := v.Object()
		require.Nil(t, err)
		t.Log(o)

		login, err := o.GetString("login")
		require.Nil(t, err)
		if login == "gophergala2016" {
			found = true
			break
		}
	}

	require.True(t, found)
}
