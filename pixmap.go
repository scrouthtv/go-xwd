package xwd

import "image/color"
import "errors"
import "fmt"
import "io"
import "log"

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
	pixels []uint8
}

func ReadPixmap(r io.Reader, h *XWDFileHeader, colors *XWDColorMap) (XWDPixmap, error) {
	if h.IsMapped() {
		log.Println("this is a mapped pixmap")
		return readPixmapMapped(r, h, colors)
	} else {
		log.Println("this is a raw pixmap")
		return readPixmapRaw(r, h)
	}
}

func readPixmapRaw(r io.Reader, h *XWDFileHeader) (*xwdPixmapRaw, error) {
	panic("not impl")
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

	return nil, nil
}

func readPixmapMapped(r io.Reader, h *XWDFileHeader, colors *XWDColorMap) (*xwdPixmapMapped, error) {
	pix := xwdPixmapMapped{}
	pix.header = h
	pix.colors = colors
	pix.pixels = make([]uint8, h.PixmapWidth * h.PixmapHeight)

	log.Println("Going to read", h.PixmapWidth * h.PixmapHeight * colormapKeySize, "bytes")
	buf := make([]byte, h.PixmapWidth * h.PixmapHeight * colormapKeySize)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	var i uint32 = 0
	var x, y uint32
	for y = 0; y < h.PixmapHeight; y++ {
		for x = 0; x < h.PixmapWidth; x++ {
			// TODO why is it a uint8 and not uin32 as suggested by the header??
			pix.pixels[i] = uint8(buf[i * colormapKeySize])
			i++
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
