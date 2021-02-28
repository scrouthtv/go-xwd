package xwd

import "image"
import "image/color"
import "image/color/palette"
import "log"
import "errors"
import "fmt"
import "io"
import "encoding/binary"

type XWDPixmap interface {
	At(x, y int) color.Color
}

// xwdPixmapRaw is used if the image's data is stored raw, as in each pixel
// stores it's own color.
type xwdPixmapRaw struct {
	header *XWDFileHeader
	pixels []XWDColor
}

// xwdPixmapMapped is used if the image's data is colormapped, e.g. 
// each pixel is a "pointer" to a value of the corresponding color map.
type xwdPixmapMapped struct {
	header *XWDFileHeader
	colors *XWDColorMap
	pixels []uint32
}

func ReadPixmap(r io.Reader, h *XWDFileHeader, colors *XWDColorMap) (XWDPixmap, error) {
	if h.IsMapped() {
		return readPixmapMapped(r, h, colors)
	} else {
		return readPixmapRaw(r, h)
	}
}

func readPixmapRaw(r io.Reader, h *XWDFileHeader) (*xwdPixmapRaw, error) {
	buf := make([]byte, h.BitsPerPixel)

	var x, y uint32
	for x = 0; x < h.PixmapWidth; x++ {
		for y = 0; y < h.PixmapHeight; y++ {
			_, err := r.Read(buf)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error reading %d %d: %s", x, y, err.Error()))
			}
		}
	}

	panic("not impl")
}

func readPixmapMapped(r io.Reader, h *XWDFileHeader, colors *XWDColorMap) (*xwdPixmapMapped, error) {
	pix := xwdPixmapMapped{}
	pix.header = h
	pix.colors = colors
	pix.pixels = make([]uint32, h.PixmapWidth * h.PixmapHeight)

	buf := make([]byte, h.PixmapWidth * h.PixmapHeight * 1) // color map key is a uint8 which has 1 byte
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	var i uint32
	var x, y uint32
	for y = 0; y < h.PixmapHeight; y++ {
		for x = 0; x < h.PixmapWidth; x++ {
			i = y * h.PixmapHeight + x
			// TODO use ByteOrder from the header
			pix.pixels[i] = binary.BigEndian.Uint32(buf[i * colormapKeySize : (i+1) * colormapKeySize])
		}
	}

	return &pix, nil
}

func (p *xwdPixmapRaw) At(x, y int) color.Color {
	return &p.pixels[y * int(p.header.PixmapWidth) + x]
}

func (p *xwdPixmapMapped) At(x, y int) color.Color {
	id := p.pixels[y * int(p.header.PixmapWidth) + x]
	color := p.colors.Get(int(id))
	return &color
}
