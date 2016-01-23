package github

import (
	"net/url"
	"os"
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
	resp := Client.Do(t, req)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)

	o, err := jason.NewObjectFromReader(resp.Body)
	require.Nil(t, err)
	login, err := o.GetString("login")
	require.Nil(t, err)
	Login = login
}

func TestListOrgs(t *testing.T) {
	req := Client.NewRequest(t, "GET", "/user/orgs")
	resp := Client.Do(t, req)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)

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
