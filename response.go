package gophers

import (
	"io/ioutil"
	"net/http"

	"github.com/gophergala2016/gophers/jsons"
)

// Response represents HTTP response.
type Response struct {
	*http.Response
}

// JSON returns reponse body as JSON structure.
// In case of error if fails test.
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
