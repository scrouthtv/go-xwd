package xwd

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	_ "embed"
)

//go:embed 8colors.xwd
var xwd8colors []byte

//go:embed 8colors.png
var png8colors []byte

func TestHeader(t *testing.T) {
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

	t.Log("\n" + imageToString(p))

	pngImage, err := png.Decode(bytes.NewReader(png8colors))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + imageToString(pngImage))

	for x := 0; x < 4; x++ {
		for y := 0; y < 2; y++ {
			if !ColorEqual(p.At(x, y), pngImage.At(x, y)) {
				ir, ig, ib, ia := p.At(x, y).RGBA()
				sr, sg, sb, sa := pngImage.At(x, y).RGBA()
				t.Errorf("Colors @ %d/%d differ: should be %d, %d, %d, %d; is %d, %d, %d, %d",
					x, y, sr, sg, sb, sa, ir, ig, ib, ia)
			}
		}
	}

	if !t.Failed() {
		t.Log("All colors equal")
	}

	dump, err := os.Create("dump.png")
	if err != nil {
		t.Fatal(err)
	}

	err = png.Encode(dump, p)
	if err != nil {
		t.Fatal(err)
	}
}
