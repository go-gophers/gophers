package github

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gophergala2016/gophers/json"
)

func TestCreateDestroyRepo(t *testing.T) {
	t.Parallel()
	if Login == "" {
		TestGetUser(t)
	}

	// create repo
	repo := TestPrefix + Faker.UserName()
	v := Client.Post(t, "/user/repos", JSON(`{"name": %q}`, repo).Reader(), 201).JSON(t)
	assert.Equal(t, JSON(`{"name": %q, "full_name": %q}`, repo, Login+"/"+repo), v.KeepFields("name", "full_name"))
	assert.Equal(t, JSON(`{"login": %q}`, Login), v.Get("/owner").KeepFields("login"))

	// try to create repo with the same name again
	v = Client.Post(t, "/user/repos", JSON(`{"name": %q}`, repo).Reader(), 422).JSON(t)
	assert.Equal(t, JSON(`{"message": "Validation Failed"}`), v.KeepFields("message"))
	assert.Equal(t, JSON(`{"code": "custom", "field": "name"}`), v.Get("/errors/0").KeepFields("code", "field"))

	// destroy repo
	Client.Delete(t, "/repos/"+Login+"/"+repo, 204)
}
