package github

import (
	"fmt"
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
	resp := Client.Do(t, req)
	defer resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode)

	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)
	resp = Client.Do(t, req)
	defer resp.Body.Close()
	require.Equal(t, 204, resp.StatusCode)
}
