package xwd

import (
	"bytes"
	"image/png"
	"testing"

	_ "embed"
)

//go:embed 8colors.xwd
var xwd8colors []byte

//go:embed 8colors.png
var png8colors []byte

func TestMappedImage(t *testing.T) {
	rdr := bytes.NewReader(xwd8colors)

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

	pngImage, err := png.Decode(bytes.NewReader(png8colors))
	if err != nil {
		t.Fatal(err)
	}

	/*t.Log("\n" + imageToString(p))

	t.Log("\n" + imageToString(pngImage))*/

	if !ImageEqual(p, pngImage, 0) {
		t.Fatal("images differ")
	}

	/*if !t.Failed() {
		t.Log("All colors equal")
	}

	dump, err := os.Create("dump.png")
	if err != nil {
		t.Fatal(err)
	}

	err = png.Encode(dump, p)
	if err != nil {
		t.Fatal(err)
	}*/
}
