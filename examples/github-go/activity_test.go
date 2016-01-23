package github

import (
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BaseURL = "https://api.github.com"
)

var (
	Token = os.Getenv("GOPHERS_GITHUB_TOKEN")
)

func makeRequest(t *testing.T, method, path string) *http.Request {
	req, err := http.NewRequest(method, BaseURL+path, nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "token "+Token)
	return req
}

func TestListOrgs(t *testing.T) {
	req := makeRequest(t, "GET", "/user/orgs")
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	b, err := httputil.DumpResponse(resp, true)
	assert.Nil(t, err)
	t.Logf("%s", b)
}
