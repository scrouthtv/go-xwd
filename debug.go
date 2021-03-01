package xwd

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
)

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
	case BigEndian:
		return "big endian (0)"
	case LittleEndian:
		return "little endian (1)"
	case Invalid:
		return "invalid"
	default:
		return "invalid (" + strconv.Itoa(int(o)) + ")"
	}
}

/*func hecateHex(p []byte) string {
	var out strings.Builder

	for _, b := range p {
		fmt.Fprintf(&out, "%02x ", b)
	}

	return out.String()
}*/

func imageToString(i image.Image, xscale int, yscale int) string {
	var out strings.Builder

	minx, miny := i.Bounds().Min.X, i.Bounds().Min.Y
	maxx, maxy := i.Bounds().Max.X, i.Bounds().Max.Y

	for y := miny; y < maxy; y += yscale {
		for x := minx; x < maxx; x += xscale {
			r, g, b, _ := i.At(x, y).RGBA()
			sr, sg, sb := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			fmt.Fprintf(&out, "\x1b[48;2;%d;%d;%dm  ", sr, sg, sb)
		}
		out.WriteString("\x1b[49m\n")
	}

	return out.String()
}

func ColorEqual(a color.Color, b color.Color) bool {
	return ColorSimilar(a, b, 0xffff)
}

func ColorSimilar(a color.Color, b color.Color, mask uint32) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()

	return ar & mask == br & mask && ag & mask == bg & mask && ab & mask == bb & mask && aa & mask == ba & mask
}

// ImageEqual tests if two images are similar:
// If the images don't have the same size, false is returned.
// Else every pixel is checked if it is similar, e.g. if all its colors are equal
// discarding the prox least bytes (since images always differ a bit).
func ImageEqual(a image.Image, b image.Image, prox int) bool {
	ax0, ay0 := a.Bounds().Min.X, a.Bounds().Min.Y
	bx0, by0 := b.Bounds().Min.X, b.Bounds().Min.Y
	ax1, ay1 := a.Bounds().Max.X, a.Bounds().Max.Y
	bx1, by1 := b.Bounds().Max.X, b.Bounds().Max.Y

	awidth := ax1 - ax0
	aheight := ay1 - ay0
	bwidth := bx1 - bx0
	bheight := by1 - by0

	if awidth != bwidth || aheight != bheight {
		debugf("Sizes differ: %dx%d / %dx%d", awidth, aheight, bwidth, bheight)
		return false
	}

	var mask uint16 = 0xffff << prox

	bxoffset := bx0 - ax0
	byoffset := by0 - ay0

	var ac, bc color.Color

	var ok bool = true

	for x := ax0; x < ax1; x++ {
		for y := ay0; y < ay1; y++ {
			ac = a.At(x, y)
			bc = b.At(x + bxoffset, y + byoffset)

			if !ColorSimilar(ac, bc, uint32(mask)) {
				ar, ag, ab, aa := ac.RGBA()
				br, bg, bb, ba := bc.RGBA()
				debugf("Colors @ %d/%d differ: %d, %d, %d, %d / %d, %d, %d, %d", x, y, ar, ag, ab, aa, br, bg, bb, ba)
				ok = false
			}
		}
	}

	return ok
}
