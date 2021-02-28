package xwd

import (
	"io"
	"log"
	"image/color"
	"errors"
	"encoding/binary"
)

// XWDColor is a color in the xwd image.
type XWDColor struct {
	Pixel uint32
	Red         uint16
	Green       uint16
	Blue        uint16
	Flags       uint8
	Padding     uint8
}

const colorSize = 4+2+2+2+1+1

var opaqAlpha uint32 = uint32(color.Opaque.A) | uint32(color.Opaque.A) << 16

// RGBA implements the image/color.Color.RGBA() method
func (c *XWDColor) RGBA() (r, g, b, a uint32) {
	return uint32(c.Red) << 16, uint32(c.Green) << 16, uint32(c.Blue) << 16, 0xffffffff
}

type XWDColorMap []XWDColor

const colormapKeySize = 1 // uint8, this is the size in the pixmap

func (c *XWDColorMap) Get(i int) XWDColor {
	return (*c)[i]
}

// ReadColorMap reads the colormap from an xwd image, provided that the header has already been read
// and discarded from the reader.
// It returns the color map or any encountered error.
func ReadColorMap(r io.Reader, h *XWDFileHeader) (XWDColorMap, error) {
	var m XWDColorMap = make([]XWDColor, h.NumberOfColors)

	r.Read(make([]byte, 0))

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
		m[i] = XWDColor{
			binary.BigEndian.Uint32(buf[0:4]),// << 8 seems to be wrong
			binary.BigEndian.Uint16(buf[4:6]),
			binary.BigEndian.Uint16(buf[6:8]),
			binary.BigEndian.Uint16(buf[8:10]),
			uint8(buf[10]),
			uint8(buf[11]),
		}
		log.Printf("Read %sas %s", hecateHex(buf), m[i].String())
	}

	return m, nil
}
