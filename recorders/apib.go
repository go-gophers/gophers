package recorders

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
	"text/template"
)

var (
	apibRequestTemplate = template.Must(template.New("apibRequest").Parse(strings.TrimSpace(`
+ Request ({{ .ContentType }})

{{ .Body }}
`)))

	apibResponseTemplate = template.Must(template.New("apibResponse").Parse(strings.TrimSpace(`
+ Response {{ .StatusCode }}

    + Headers

{{ .Headers }}

    + Body

{{ .Body }}
`)))
)

// APIB records request and response to writers in API Blueprint format.
type APIB struct{}

func (r *APIB) RecordRequest(req *http.Request, status, headers, body []byte, wc io.WriteCloser) (err error) {
	indent := strings.Repeat(" ", 13)

	// indent body
	var bodyS []string
	s := bufio.NewScanner(bytes.NewReader(body))
	for s.Scan() {
		bodyS = append(bodyS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	err = apibRequestTemplate.Execute(wc, map[string]interface{}{
		"ContentType": req.Header.Get("Content-Type"),
		"Body":        strings.Join(bodyS, "\n"),
	})
	if err == nil {
		err = wc.Close()
	}
	return
}

func (r *APIB) RecordResponse(resp *http.Response, status, headers, body []byte, wc io.WriteCloser) (err error) {
	indent := strings.Repeat(" ", 12)

	// indent headers
	var headersS []string
	s := bufio.NewScanner(bytes.NewReader(headers))
	for s.Scan() {
		headersS = append(headersS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	// indent body
	var bodyS []string
	s = bufio.NewScanner(bytes.NewReader(body))
	for s.Scan() {
		bodyS = append(bodyS, indent+s.Text())
	}
	if err = s.Err(); err != nil {
		return
	}

	err = apibResponseTemplate.Execute(wc, map[string]interface{}{
		"StatusCode": resp.StatusCode,
		"Headers":    strings.Join(headersS, "\n"),
		"Body":       strings.Join(bodyS, "\n"),
	})
	if err == nil {
		err = wc.Close()
	}
	return
}

// check interface
var _ Interface = new(APIB)
