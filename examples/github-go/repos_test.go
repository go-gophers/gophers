package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDestroyRepo(t *testing.T) {
	if Login == "" {
		TestGetUser(t)
	}

	repo := TestPrefix + Faker.UserName()
	req := Client.NewRequest(t, "POST", "/user/repos")
	req.Body = ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"name": %q}`, repo)))
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

	// TODO check status code

	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)
	b, err = httputil.DumpRequestOut(req, true)
	require.Nil(t, err)
	t.Logf("Request:\n%s", b)

	resp, err = http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	require.Nil(t, err)
	b, err = httputil.DumpResponse(resp, true)
	require.Nil(t, err)
	t.Logf("Response:\n%s", b)

	// TODO check status code
}
