package github

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gophergala2016/gophers/json"
)

func TestGetCurrentUser(t *testing.T) {
	j := Client.Get(t, "/user", 200).JSON(t)
	Login = j.KeepFields("login").(JSONObject)["login"].(string)
	require.NotEmpty(t, Login)
}
