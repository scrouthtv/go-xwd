package xwd

import "testing"
import "bytes"
import "strings"
import "fmt"

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

	p, err := ReadPixmap(rdr, hdr, &colors)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.At(3, 1))

	/*t.Logf("\nPixmap (%d bytes):\n", len(*p))
	t.Logf("%x\n", *p)*/

	fmt.Println(p.At(0, 0).(*XWDColor).String())
	var x, y uint32
	var out strings.Builder
	for y = 0; y < hdr.PixmapHeight; y++ {
		for x = 0; x < hdr.PixmapWidth; x++ {
			r, g, b, _ := p.At(int(x), int(y)).RGBA()
			sr, sg, sb := uint8(r >> 24), uint8(g >> 24), uint8(b >> 24)
			fmt.Fprintf(&out, "\x1b[48;2;%d;%d;%dm  ", sr, sg, sb)
			if x == 1 && y == 0 {
				fmt.Println(sr, sg, sb)
			}
		}
		out.WriteString("\x1b[49m\n")
	}

	//if false {
	fmt.Println(out.String())
//}
}
