package xwd

import (
	"image"
	"io"
	"log"
)

// Matches any header size (uint32), fileversion (uint32) == 7.
const xwdHeader = "????0007"

// DoDebug indicates whether `go-xwd` should print any debugging messages.
// If enabled, debugging messages are print using log.Print().
var DoDebug = false

func init() {
	image.RegisterFormat("xwd", xwdHeader, Decode, DecodeConfig)
}

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

// DecodeConfig tries to read an xwd image and returns it's config.
func DecodeConfig(r io.Reader) (image.Config, error) {
	hdr, err := ReadHeader(r)
	if err != nil {
		return image.Config{}, err
	}
	return hdr.Config(), nil
}

func debugf(format string, a ...interface{}) {
	if DoDebug {
		log.Printf(format, a...)
	}
}
