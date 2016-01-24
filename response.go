package gophers

import (
	"io/ioutil"
	"net/http"

	"github.com/gophergala2016/gophers/jsons"
)

type Response struct {
	*http.Response
}

func (r *Response) JSON(t TestingTB) (j jsons.Struct) {
	defer func() {
		if p := recover(); p != nil {
			j = nil
			t.Fatalf("panic: %v", p)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("can't read body: %s", err)
	}

	j = jsons.Parse(string(b))
	return
}
