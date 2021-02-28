package xwd

import "testing"
import "bytes"
import _ "embed"

//go:embed 500colors.xwd
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

	t.Logf("Color map (%d entries):", len(colors))
	for _, c := range colors {
		t.Log(c.String())
	}

	p, err := ReadPixmap(rdr, hdr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.At(0, 0))

	/*t.Logf("\nPixmap (%d bytes):\n", len(*p))
	t.Logf("%x\n", *p)*/
}
