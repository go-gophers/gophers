package glua

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yuin/gopher-lua"

	"github.com/go-gophers/gophers"
)

type client struct {
	c *http.Client
}

func NewClient(c *http.Client) *client {
	return &client{
		c: c,
	}
}

func (c *client) Loader(state *lua.LState) int {
	module := state.SetFuncs(state.NewTable(), map[string]lua.LGFunction{
		"get": c.get,
	})
	requestRegister(module, state)
	responseRegister(module, state)
	state.Push(module)
	return 1
}

func (c *client) get(state *lua.LState) int {
	return c.doAndPush(state, "GET", state.ToString(1), state.ToTable(2))
}

func (c *client) doAndPush(state *lua.LState, method string, url string, options *lua.LTable) int {
	response, err := c.do(state, method, url, options)

	if err != nil {
		state.Push(lua.LNil)
		state.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	state.Push(response)
	return 1
}

func (c *client) do(state *lua.LState, method string, url string, options *lua.LTable) (*lua.LUserData, error) {
	// TODO use Client.NewRequest
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req := &gophers.Request{Request: r}

	var body []byte
	resp, err := c.c.Do(req.Request)
	if resp != nil {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return responseNew(&gophers.Response{Response: resp}, body, state), nil
}
