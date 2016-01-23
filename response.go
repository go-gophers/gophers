package gophers

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/gophergala2016/gophers/json"
)

type Response struct {
	*http.Response
}

func (r *Response) JSON(t testing.TB) (j JSONStruct) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
			j = nil
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	j = JSON(string(b))
	return
}
