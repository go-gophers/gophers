package github

import (
	"fmt"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDestroyRepo(t *testing.T) {
	if Login == "" {
		TestGetUser(t)
	}

	repo := TestPrefix + Faker.UserName()
	req := Client.NewRequest(t, "POST", "/user/repos")
	req.SetBodyString(fmt.Sprintf(`{"name": %q}`, repo))

	b, err := httputil.DumpRequestOut(req.Request, true)
	require.Nil(t, err)
	t.Logf("Request:\n%s", b)

	resp := Client.Do(t, req)
	defer resp.Body.Close()

	b, err = httputil.DumpResponse(resp, true)
	require.Nil(t, err)
	t.Logf("Response:\n%s", b)

	// TODO check status code

	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)

	b, err = httputil.DumpRequestOut(req.Request, true)
	require.Nil(t, err)
	t.Logf("Request:\n%s", b)

	resp = Client.Do(t, req)
	defer resp.Body.Close()

	b, err = httputil.DumpResponse(resp, true)
	require.Nil(t, err)
	t.Logf("Response:\n%s", b)

	// TODO check status code
}
