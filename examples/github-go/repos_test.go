package github

import (
	"testing"

	"github.com/manveru/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-gophers/gophers/jsons"
)

func createRepo(t *testing.T, record bool) string {
	// Create new Faker instance since it's not thread-safe
	// https://github.com/manveru/faker/issues/6

	faker, err := faker.New("en")
	require.Nil(t, err)

	// create repo
	repo := TestPrefix + faker.UserName()
	req := Client.NewRequest(t, "POST", "/user/repos", jsons.Parse(`{"name": %q}`, repo))
	if record {
		req.EnableRecording("repo_create.apib")
	}
	j := Client.Do(t, req, 201).JSON(t)
	assert.Equal(t, jsons.Parse(`{"name": %q, "full_name": %q}`, repo, Login+"/"+repo), j.KeepFields("name", "full_name"))
	assert.Equal(t, jsons.Parse(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))
	return repo
}

func destroyRepo(t *testing.T, repo string) {
	Client.Delete(t, "/repos/"+Login+"/"+repo, 204)
}

func TestRepoCreateDestroy(t *testing.T) {
	t.Parallel()

	repo := createRepo(t, true)
	defer destroyRepo(t, repo)

	// check repo exists
	j := Client.Get(t, "/repos/"+Login+"/"+repo, 200).JSON(t)
	assert.Equal(t, jsons.Parse(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))

	// try to create repo with the same name again
	req := Client.NewRequest(t, "POST", "/user/repos", jsons.Parse(`{"name": %q}`, repo)).EnableRecording("repo_create_exist.apib")
	j = Client.Do(t, req, 422).JSON(t)
	assert.Equal(t, jsons.Parse(`{"message": "Validation Failed"}`), j.KeepFields("message"))
	assert.Equal(t, jsons.Parse(`{"code": "custom", "field": "name"}`), j.Get("/errors/0").KeepFields("code", "field"))
}

func TestRepoList(t *testing.T) {
	t.Parallel()

	repo := createRepo(t, false)
	defer destroyRepo(t, repo)

	j := Client.Get(t, "/user/repos?visibility=public&affiliation=owner&sort=created", 200).JSON(t)

	var found bool
	v := j.KeepFields("name")
	expect := jsons.Parse(`{"name": %q}`, repo).String()
	for _, e := range v.(jsons.Array) {
		if jsons.Cast(e).String() == expect {
			found = true
			break
		}
	}

	assert.True(t, found, "created repository not found in list")
}
