package xwd

import "testing"
import "bytes"
import _ "embed"

//go:embed good.xwd
var xwd8colors []byte

func TestHeader(t *testing.T) {
	rdr := bytes.NewReader(xwd8colors)
	hdr, err := ReadHeader(rdr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hdr)
}
