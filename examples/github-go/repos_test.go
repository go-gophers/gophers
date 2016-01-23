package github

import (
	"fmt"
	"testing"
)

func TestCreateDestroyRepo(t *testing.T) {
	t.Parallel()
	if Login == "" {
		TestGetUser(t)
	}

	repo := TestPrefix + Faker.UserName()
	req := Client.NewRequest(t, "POST", "/user/repos")
	req.SetBodyString(fmt.Sprintf(`{"name": %q}`, repo))
	resp := Client.Do(t, req, 201)

	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)
	resp = Client.Do(t, req, 204)
	_ = resp
}
