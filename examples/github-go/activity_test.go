package github

import (
	"net/url"
	"os"
	"testing"

	"github.com/manveru/faker"
	"github.com/stretchr/testify/require"

	"github.com/gophergala2016/gophers"
	. "github.com/gophergala2016/gophers/json"
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
	v := Client.Do(t, req, 200).JSON(t).KeepFields("login")
	Login = v.(JSONObject)["login"].(string)
}

func TestListOrgs(t *testing.T) {
	t.Parallel()

	req := Client.NewRequest(t, "GET", "/user/orgs")
	v := Client.Do(t, req, 200).JSON(t).KeepFields("login")

	var found bool
	expect := JSON(`{"login": "gophergala2016"}`).String()
	for _, e := range v.(JSONArray) {
		if AsJSON(e).String() == expect {
			found = true
			break
		}
	}

	require.True(t, found)
}
