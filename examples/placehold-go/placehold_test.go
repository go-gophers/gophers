package placehold

import (
	"bytes"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func basicTest(t testing.TB) {
	req := Client.NewRequest(t, "GET", "/350x150", nil)
	req.EnableRecording("placehold.apib")
	res := Client.Do(t, req, 200)
	img, format, err := image.Decode(bytes.NewReader(res.Body))
	assert.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, "(350,150)", img.Bounds().Max.String())

	res = Client.Get(t, "/350x150.jpeg", 200)
	img, format, err = image.Decode(bytes.NewReader(res.Body))
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", format)
	assert.Equal(t, "(350,150)", img.Bounds().Max.String())
}

func TestBasic(t *testing.T) {
	basicTest(t)
}

func BenchmarkBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basicTest(b)
	}
}
