package github

import (
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophergala2016/gophers"
)

var (
	Client = gophers.NewClient(*gophers.MustParseURL("https://api.github.com/?access_token=" + os.Getenv("GOPHERS_GITHUB_TOKEN")))
)

func TestListOrgs(t *testing.T) {
	req := Client.NewRequest(t, "GET", "/user/orgs")
	b, err := httputil.DumpRequestOut(req, true)
	assert.Nil(t, err)
	t.Logf("Request:\n%s", b)

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	b, err = httputil.DumpResponse(resp, true)
	assert.Nil(t, err)
	t.Logf("Response:\n%s", b)
}
