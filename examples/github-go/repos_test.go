package github

import (
	"fmt"
	"testing"

	"github.com/antonholmquist/jason"
	"github.com/stretchr/testify/assert"
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
	o, err := jason.NewObjectFromReader(resp.Body)
	assert.Nil(t, err)
	login, err := o.GetString("owner", "login")
	assert.Nil(t, err)
	assert.Equal(t, Login, login)
	name, err := o.GetString("name")
	assert.Nil(t, err)
	assert.Equal(t, repo, name)
	name, err = o.GetString("full_name")
	assert.Nil(t, err)
	assert.Equal(t, Login+"/"+repo, name)

	// try to create repo with the same name again
	req = Client.NewRequest(t, "POST", "/user/repos")
	req.SetBodyString(fmt.Sprintf(`{"name": %q}`, repo))
	resp = Client.Do(t, req, 422)

	// check response
	o, err = jason.NewObjectFromReader(resp.Body)
	assert.Nil(t, err)
	message, err := o.GetString("message")
	assert.Nil(t, err)
	assert.Equal(t, "Validation Failed", message)
	errors, err := o.GetObjectArray("errors")
	assert.Nil(t, err)
	assert.Len(t, errors, 1)
	expected := `{
		"resource": "Repository",
		"code": "custom",
		"field": "name",
		"message": "name already exists on this account"
	}`
	assert.JSONEq(t, expected, errors[0].String())

	// destroy repo
	req = Client.NewRequest(t, "DELETE", "/repos/"+Login+"/"+repo)
	Client.Do(t, req, 204)
}
