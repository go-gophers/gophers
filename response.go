package gophers

import (
	"io/ioutil"
	"net/http"

	. "github.com/gophergala2016/gophers/json"
)

type Response struct {
	*http.Response
}

func (r *Response) JSON(t TestingTB) (j JSONStruct) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("panic: %v", p)
			j = nil
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("can't read body: %s", err)
		return
	}

	j = JSON(string(b))
	return
}
