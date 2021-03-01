package xwd

import "strings"
import "fmt"
import "strconv"
import "image"
import "image/color"

// String creates a textual representation of this header.
// It is comparable to the output of xwud -dumpheaders.
func (h *FileHeader) String() string {
	var out strings.Builder

	fmt.Fprintf(&out, "window name:        %s\n", h.WindowName)
	fmt.Fprintf(&out, "sizeof(XWDheader):  %d\n", xwdHeaderSize)
	fmt.Fprintf(&out, "header size:        %d\n", h.HeaderSize)
	fmt.Fprintf(&out, "file version:       %d\n", h.FileVersion)
	fmt.Fprintf(&out, "pixmap format:      %d\n", h.PixmapFormat)
	fmt.Fprintf(&out, "pixmap depth:       %d\n", h.PixmapDepth)
	fmt.Fprintf(&out, "pixmap width:       %d\n", h.PixmapWidth)
	fmt.Fprintf(&out, "pixmap height:      %d\n", h.PixmapHeight)
	fmt.Fprintf(&out, "x offset:           %d\n", h.XOffset)
	fmt.Fprintf(&out, "byte order:         %d\n", h.ByteOrder)
	fmt.Fprintf(&out, "bitmap unit:        %d\n", h.BitmapUnit)
	fmt.Fprintf(&out, "bitmap bit order:   %d\n", h.BitmapBitOrder)
	fmt.Fprintf(&out, "bitmap pad:         %d\n", h.BitmapPad)
	fmt.Fprintf(&out, "bits per pixel:     %d\n", h.BitsPerPixel)
	fmt.Fprintf(&out, "bytes per line:     %d\n", h.BytesPerLine)
	fmt.Fprintf(&out, "visual class:       %d\n", h.VisualClass)
	fmt.Fprintf(&out, "red mask:           %d\n", h.RedMask)
	fmt.Fprintf(&out, "green mask:         %d\n", h.GreenMask)
	fmt.Fprintf(&out, "blue mask:          %d\n", h.BlueMask)
	fmt.Fprintf(&out, "bits per rgb:       %d\n", h.BitsPerRgb)
	fmt.Fprintf(&out, "colormap entries:   %d\n", h.ColorMapEntries)
	fmt.Fprintf(&out, "num colors:         %d\n", h.NumberOfColors)
	fmt.Fprintf(&out, "window width:       %d\n", h.WindowWidth)
	fmt.Fprintf(&out, "window height:      %d\n", h.WindowHeight)
	fmt.Fprintf(&out, "window x:           %d\n", h.WindowX)
	fmt.Fprintf(&out, "window y:           %d\n", h.WindowY)
	fmt.Fprintf(&out, "border width:       %d\n", h.WindowBorderWidth)

	return out.String()
}

func (c *Color) String() string {
	return fmt.Sprintf("Color %d: %d/%d/%d (%d)", c.Pixel, c.Red, c.Green, c.Blue, c.Flags)
}

func (o Order) String() string {
	switch o {
	case BigEndian: return "big endian (0)"
	case LittleEndian: return "little endian (1)"
	default: return "invalid (" + strconv.Itoa(int(o)) + ")"
	}
}

/*func hecateHex(p []byte) string {
	var out strings.Builder

	for _, b := range p {
		fmt.Fprintf(&out, "%02x ", b)
	}

	return out.String()
}*/

func imageToString(i image.Image) string {
	var out strings.Builder

	minx, miny := i.Bounds().Min.X, i.Bounds().Min.Y
	maxx, maxy := i.Bounds().Max.X, i.Bounds().Max.Y

	for y := miny; y < maxy; y++ {
		for x := minx; x < maxx; x++ {
			r, g, b, _ := i.At(x, y).RGBA()
			sr, sg, sb := uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
			fmt.Fprintf(&out, "\x1b[48;2;%d;%d;%dm  ", sr, sg, sb)
		}
		out.WriteString("\x1b[49m\n")
	}

	return out.String()
}

func ColorEqual(a color.Color, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
