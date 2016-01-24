package github

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gophergala2016/gophers/json"
)

func createRepo(t *testing.T) string {
	// create repo
	repo := TestPrefix + Faker.UserName()
	j := Client.Post(t, "/user/repos", JSON(`{"name": %q}`, repo).Reader(), 201).JSON(t)
	assert.Equal(t, JSON(`{"name": %q, "full_name": %q}`, repo, Login+"/"+repo), j.KeepFields("name", "full_name"))
	assert.Equal(t, JSON(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))
	return repo
}

func destroyRepo(t *testing.T, repo string) {
	Client.Delete(t, "/repos/"+Login+"/"+repo, 204)
}

func TestRepoCreateDestroy(t *testing.T) {
	t.Parallel()

	repo := createRepo(t)
	defer destroyRepo(t, repo)

	// check repo exists
	j := Client.Get(t, "/repos/"+Login+"/"+repo, 200).JSON(t)
	assert.Equal(t, JSON(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))

	// try to create repo with the same name again
	j = Client.Post(t, "/user/repos", JSON(`{"name": %q}`, repo).Reader(), 422).JSON(t)
	assert.Equal(t, JSON(`{"message": "Validation Failed"}`), j.KeepFields("message"))
	assert.Equal(t, JSON(`{"code": "custom", "field": "name"}`), j.Get("/errors/0").KeepFields("code", "field"))
}

func TestRepoList(t *testing.T) {
	t.Parallel()

	repo := createRepo(t)
	defer destroyRepo(t, repo)

	j := Client.Get(t, "/user/repos?visibility=public&affiliation=owner&sort=created", 200).JSON(t)

	var found bool
	v := j.KeepFields("name")
	expect := JSON(`{"name": %q}`, repo).String()
	for _, e := range v.(JSONArray) {
		if AsJSON(e).String() == expect {
			found = true
			break
		}
	}

	assert.True(t, found, "created repository not found in list")
}
