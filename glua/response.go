package glua

import (
	"github.com/yuin/gopher-lua"

	"github.com/go-gophers/gophers"
)

const responseType = "gophers.response"

type response struct {
	*gophers.Response
	b lua.LString
}

func (resp *response) code(state *lua.LState) int {
	state.Push(lua.LNumber(resp.StatusCode))
	return 1
}

func (resp *response) status(state *lua.LState) int {
	state.Push(lua.LString(resp.Status))
	return 1
}

func (resp *response) body(state *lua.LState) int {
	state.Push(lua.LString(resp.b))
	return 1
}

func responseRegister(module *lua.LTable, state *lua.LState) {
	mt := state.NewTypeMetatable(responseType)
	state.SetField(mt, "__index", state.NewFunction(responseIndex))
	state.SetField(module, "response", mt)
}

func responseNew(resp *gophers.Response, body []byte, state *lua.LState) *lua.LUserData {
	ud := state.NewUserData()
	ud.Value = &response{
		Response: resp,
		b:        lua.LString(body),
	}
	state.SetMetatable(ud, state.GetTypeMetatable(responseType))
	return ud
}

func responseCheck(state *lua.LState) *response {
	ud := state.CheckUserData(1)
	if v, ok := ud.Value.(*response); ok {
		return v
	}
	state.ArgError(1, responseType+" expected")
	return nil
}

func responseIndex(state *lua.LState) int {
	resp := responseCheck(state)

	switch state.CheckString(2) {
	case "code":
		return resp.code(state)
	case "status":
		return resp.status(state)
	case "body":
		return resp.body(state)
	}

	return 0
}
