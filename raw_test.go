package xwd

import (
	"bytes"
	"image/png"
	"testing"

	_ "embed"
)

//go:embed 500colors.xwd
var xwd500colors []byte

//go:embed 500colors.png
var png500colors []byte

func TestRawImage(t *testing.T) {
	rdr := bytes.NewReader(xwd500colors)

	hdr, err := ReadHeader(rdr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hdr)

	colors, err := ReadColorMap(rdr, hdr)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Color map (%d entries):", len(colors))

	for _, c := range colors {
		t.Log(c.String())
	}

	p, err := ReadPixmap(rdr, hdr, &colors)
	if err != nil {
		t.Fatal(err)
	}

	pngImage, err := png.Decode(bytes.NewReader(png500colors))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + imageToString(p, 1, 1))

	t.Log("\n" + imageToString(pngImage, 1, 1))

	if !ImageEqual(p, pngImage, 8) { // have to approximate by 8 because png has a greater bit depth
		t.Fatal("Images differ")
	}
}
