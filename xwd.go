package xwd

import (
	"image"
	"image/color"
	"io"
)

// XWDImage groups together an xwd header and
// a paletted image.
// It's pointer type implements all image.Image functionality.
type XWDImage struct {
	header XWDFileHeader
	image *image.Paletted
}

func (xwd *XWDImage) Header() XWDFileHeader {
	return xwd.header
}

func (xwd *XWDImage) At(x, y int) color.Color {
	return xwd.image.At(x, y)
}

func (xwd *XWDImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(xwd.header.PixmapWidth), int(xwd.header.PixmapHeight))
}

func (xwd *XWDImage) ColorModel() color.Model {
	return color.RGBAModel
}

// Decode reads a XWD image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	//xwd := XWDImage{}

	panic("not impl")

	return nil, nil
}
