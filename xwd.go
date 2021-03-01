package xwd

import (
	"image"
	"io"
)

// Decode reads a XWD image from r and returns it as an image.Image.
// Reading happens in three steps:
// 1. Read the header (including the window name).
// 2. Read the colormap.
// 3. Read the buffer / pixmap.
func Decode(r io.Reader) (image.Image, error) {
	hdr, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}

	colors, err := ReadColorMap(r, hdr)
	if err != nil {
		return nil, err
	}

	pix, err := ReadPixmap(r, hdr, &colors)
	if err != nil {
		return nil, err
	}

	return pix, err
}
