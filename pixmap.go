package xwd

import "image/color"
import "errors"
import "fmt"
import "io"
import "log"
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
	if h.BitsPerPixel != 24 {
		return nil, errors.New("only 24 bpp supporetd")
	}

	buf := make([]byte, 4)
	var cu uint32
	var cl XWDColor

	var pixmap xwdPixmapRaw = xwdPixmapRaw{h, make([]XWDColor, h.PixmapWidth * h.PixmapHeight)}

	var rs, gs, bs int = shiftwidth(h.RedMask), shiftwidth(h.GreenMask), shiftwidth(h.BlueMask)

	var i uint32 = 0
	var x, y uint32
	for y = 0; y < h.PixmapHeight; y++ {
		for x = 0; x < h.PixmapWidth; x++ {
			_, err := r.Read(buf[1:4])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error reading %d %d: %s", x, y, err.Error()))
			}
			cu = binary.BigEndian.Uint32(buf)
			cl = XWDColor{
				Pixel: i, Flags: 7, Padding: 0,
				Red: uint16(((cu & h.RedMask) >> rs) << 8),
				Green: uint16(((cu & h.GreenMask) >> gs) << 8),
				Blue: uint16(((cu & h.BlueMask) >> bs) << 8),
			}
			pixmap.pixels[i] = cl
			if y == 0 && x > 20 {
				fmt.Println(x, "/", y, ": Read", hecateHex(buf[1:4]))
			}
			if y == 1 && x < 5 {
				fmt.Println(x, "/", y, ": Read", hecateHex(buf[1:4]))
				fmt.Println("set", x, y, "(", i, ")")
				fmt.Println(cl.String())
			}
			i++
		}
		r.Read(buf[1:3]) // discard two bytes
	}

	return &pixmap, nil
}

func shiftwidth(mask uint32) int {
	if mask == 0 {
		return 0
	}

	for i := 0; i < 32; i++ {
		if mask & 0b1 != 0 {
			return i
		}
		mask = mask >> 1
	}

	return 0
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
