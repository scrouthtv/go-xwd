package xwd

import (
	"io"
	"image/color"
	"errors"
	"encoding/binary"
)

// Color is a color in the xwd image.
type Color struct {
	Pixel uint32
	Red         uint16
	Green       uint16
	Blue        uint16
	Flags       uint8
	Padding     uint8
}

const colorSize = 4+2+2+2+1+1

// RGBA implements the image/color.Color.RGBA() method.
// It returns rgb values in between 0 and 0xffff and an alpha value of 0xffff.
func (c *Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.Red), uint32(c.Green), uint32(c.Blue), uint32(color.Opaque.A)
}

type ColorMap []Color

const colormapKeySize = 1 // uint8, this is the size in the pixmap, not in the colormap

func (c *ColorMap) Get(i int) Color {
	return (*c)[i]
}

// ReadColorMap reads the colormap from an xwd image, provided that the header has already been read
// and discarded from the reader.
// It returns the color map or any encountered error.
func ReadColorMap(r io.Reader, h *FileHeader) (ColorMap, error) {
	var m ColorMap = make([]Color, h.NumberOfColors)

	// Use NumOfColors instead of ColorMapEntries: https://gitlab.freedesktop.org/xorg/app/xwd/-/blob/master/xwd.c#L489
	var i uint32
	buf := make([]byte, colorSize)
	var n int
	var err error
	for i = 0; i < h.NumberOfColors; i++ {
		n, err = r.Read(buf)
		if err != nil {
			return nil, err
		}
		if n != colorSize {
			return nil, errors.New("partial color read")
		}
		m[i] = Color{
			binary.BigEndian.Uint32(buf[0:4]),// << 8 seems to be wrong
			binary.BigEndian.Uint16(buf[4:6]),
			binary.BigEndian.Uint16(buf[6:8]),
			binary.BigEndian.Uint16(buf[8:10]),
			uint8(buf[10]),
			uint8(buf[11]),
		}
	}

	return m, nil
}
