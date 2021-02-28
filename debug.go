package xwd

import "strings"
import "fmt"
import "strconv"

// String creates a textual representation of this header.
// It is comparable to the output of xwud -dumpheaders.
func (h *XWDFileHeader) String() string {
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

func (o Order) String() string {
	switch o {
	case BigEndian: return "big endian (0)"
	case LittleEndian: return "little endian (1)"
	default: return "invalid (" + strconv.Itoa(int(o)) + ")"
	}
}
