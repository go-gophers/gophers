package gophers

import (
	"net/http"

	"github.com/go-gophers/gophers/jsons"
)

// Response represents HTTP response.
type Response struct {
	*http.Response
	Body []byte // filled by Client.Do
}

// JSON returns reponse body as JSON structure.
// In case of error if fails test.
func (r *Response) JSON(t TestingT) (j jsons.Struct) {
	defer func() {
		if p := recover(); p != nil {
			j = nil
			t.Fatalf("panic: %v", p)
		}
	}()

	j = jsons.ParseBytes(r.Body)
	return
}
