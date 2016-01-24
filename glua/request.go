package glua

import (
	"github.com/yuin/gopher-lua"

	"github.com/gophergala2016/gophers"
)

const requestType = "gophers.request"

type request struct {
	*gophers.Request
}

func requestRegister(module *lua.LTable, state *lua.LState) {
	mt := state.NewTypeMetatable(requestType)
	state.SetField(mt, "__index", state.NewFunction(requestIndex))
	state.SetField(module, "request", mt)
}

func requestCheck(state *lua.LState) *request {
	ud := state.CheckUserData(1)
	if v, ok := ud.Value.(*request); ok {
		return v
	}
	state.ArgError(1, requestType+" expected")
	return nil
}

func requestIndex(state *lua.LState) int {
	resp := requestCheck(state)

	// TODO
	_ = resp

	switch state.CheckString(2) {
	// case "body":
	// 	return resp.body(state)
	}

	return 0
}
