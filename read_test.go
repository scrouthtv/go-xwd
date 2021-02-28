package xwd

import "testing"
import "bytes"
import _ "embed"

//go:embed map.xwd
var xwd8colors []byte

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

	for _, c := range colors {
		t.Logf("Color %d: %d/%d/%d (%d)\n", c.Pixel, c.Red, c.Green, c.Blue, c.Flags)
	}
}
