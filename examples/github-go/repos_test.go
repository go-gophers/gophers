package github

import (
	"fmt"
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
	req := Client.NewRequest(t, "POST", "/user/repos")
	req.SetBodyString(fmt.Sprintf(`{"name": %q}`, repo))
	resp := Client.Do(t, req, 201)

	// check created repo
	v := ReadJSON(t, resp.Body)
	assert.Equal(t, JSON(`{"name": %q, "full_name": %q}`, repo, Login+"/"+repo), v.KeepFields("name", "full_name"))
	assert.Equal(t, JSON(`{"login": %q}`, Login), v.Get("/owner").KeepFields("login"))

	// try to create repo with the same name again
	req = Client.NewRequest(t, "POST", "/user/repos")
	req.SetBodyString(fmt.Sprintf(`{"name": %q}`, repo))
	resp = Client.Do(t, req, 422)

	// check response
	v = ReadJSON(t, resp.Body)
	assert.Equal(t, JSON(`{"message": "Validation Failed"}`), v.KeepFields("message"))
	assert.Equal(t, JSON(`{"code": "custom", "field": "name"}`), v.Get("/errors/0").KeepFields("code", "field"))

	// destroy repo
	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)
	Client.Do(t, req, 204)
}
