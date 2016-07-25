package gophers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBodyRepr(t *testing.T) {
	r := bodyRepr("application/json", []byte(`{"a":1}`))
	expected := []byte(`{
  "a": 1
}`)
	assert.Equal(t, expected, r, "\n%q\n%q", expected, r)

	r = bodyRepr("application/x-mpegURL", []byte(`
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=628000,CODECS="avc1.4D401E,mp4a.40.2",RESOLUTION=640x360
http://host/2695_1.m3u8
`))
	expected = []byte(`
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=628000,CODECS="avc1.4D401E,mp4a.40.2",RESOLUTION=640x360
http://host/2695_1.m3u8
`)
	assert.Equal(t, expected, r, "\n%q\n%q", expected, r)
}
