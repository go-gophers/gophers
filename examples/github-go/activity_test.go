package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/antonholmquist/jason"
	"github.com/manveru/faker"
	"github.com/stretchr/testify/require"

	"github.com/gophergala2016/gophers"
)

var (
	TestPrefix = "test-gophers-"

	Login  string
	Client *gophers.Client
	Faker  *faker.Faker
)

func init() {
	urlStr := "https://api.github.com/?access_token=" + os.Getenv("GOPHERS_GITHUB_TOKEN")
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	Client = gophers.NewClient(*u)

	Faker, err = faker.New("en")
	if err != nil {
		panic(err)
	}
}

func TestGetUser(t *testing.T) {
	req := Client.NewRequest(t, "GET", "/user")
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

	o, err := jason.NewObjectFromReader(resp.Body)
	require.Nil(t, err)
	login, err := o.GetString("login")
	require.Nil(t, err)
	Login = login
}

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

	// TODO check response status code

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
