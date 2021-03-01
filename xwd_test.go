package xwd

import "testing"
import "bytes"
import "image/png"

import _ "embed"

//go:embed small.xwd
var smallxwd []byte

//go:embed small.png
var smallpng []byte

func TestXwdImage(t *testing.T) {
	png, err := png.Decode(bytes.NewReader(smallpng))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("png:\n" + imageToString(png, 4, 4))

	xwd, err := Decode(bytes.NewReader(smallxwd))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("xwd:\n" + imageToString(xwd, 4, 4))
}
