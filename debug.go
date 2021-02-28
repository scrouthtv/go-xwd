package xwd

import "strings"
import "fmt"
import "strconv"

func (h *XWDFileHeader) String() string {
	var out strings.Builder

	fmt.Fprintf(&out, "window name:        %s", h.WindowName)
	fmt.Fprintf(&out, "sizeof(XWDheader):  %d", xwdHeaderSize)
	fmt.Fprintf(&out, "header size:        %d", h.HeaderSize)
	fmt.Fprintf(&out, "file version:       %d", h.FileVersion)
	fmt.Fprintf(&out, "pixmap format:      %d", h.PixmapFormat)
	fmt.Fprintf(&out, "pixmap depth:       %d", h.PixmapDepth)
	fmt.Fprintf(&out, "pixmap width:       %d", h.PixmapWidth)
	fmt.Fprintf(&out, "pixmap height:      %d", h.PixmapHeight)
	fmt.Fprintf(&out, "x offset:           %d", h.XOffset)
	fmt.Fprintf(&out, "byte order:         %d", h.ByteOrder)
	fmt.Fprintf(&out, "bitmap unit:        %d", h.BitmapUnit)
	fmt.Fprintf(&out, "bitmap bit order:   %d", h.BitmapBitOrder)
	fmt.Fprintf(&out, "bitmap pad:         %d", h.BitmapPad)
	fmt.Fprintf(&out, "bits per pixe:      %d", h.BitsPerPixel)
	fmt.Fprintf(&out, "visual class:       %d", h.VisualClass)
	fmt.Fprintf(&out, "red mask:           %d", h.RedMask)
	fmt.Fprintf(&out, "green mask:         %d", h.GreenMask)
	fmt.Fprintf(&out, "blue mask:          %d", h.BlueMask)
	fmt.Fprintf(&out, "bits per rgb:       %d", h.BitsPerRgb)
	fmt.Fprintf(&out, "colormap entries:   %d", h.ColorMapEntries)
	fmt.Fprintf(&out, "num colors:         %d", h.NumberOfColors)
	fmt.Fprintf(&out, "window width:       %d", h.WindowWidth)
	fmt.Fprintf(&out, "window height:      %d", h.WindowHeight)
	fmt.Fprintf(&out, "window x:           %d", h.WindowX)
	fmt.Fprintf(&out, "window y:           %d", h.WindowY)
	fmt.Fprintf(&out, "border width:       %d", h.WindowBorderWidth)

	return out.String()
}

func (o Order) String() string {
	switch o {
	case BigEndian: return "big endian (0)"
	case LittleEndian: return "little endian (1)"
	default: return "invalid (" + strconv.Itoa(int(o)) + ")"
	}
}
