package placehold

import (
	"bytes"
	"image"

	"github.com/go-gophers/gophers"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t gophers.TestingT) {
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

func TestFail(t gophers.TestingT) {
	Client.NewRequest(t, "GET", "/350x150", nil)
	t.Fatal("fatal error")
}

func TestPanic(t gophers.TestingT) {
	Client.NewRequest(t, "GET", "/350x150", nil)
	panic("PANIC!")
}
