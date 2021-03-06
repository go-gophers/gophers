package github

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-gophers/gophers/jsons"
)

func TestGetCurrentUser(t *testing.T) {
	j := Client.Get(t, "/user", 200).JSON(t)
	Login = j.KeepFields("login").(jsons.Object)["login"].(string)
	require.NotEmpty(t, Login)
}
