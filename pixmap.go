package xwd

import "image"
import "image/color"
import "errors"
import "fmt"
import "io"
import "encoding/binary"

type Pixmap interface {
	At(x, y int) color.Color
	Bounds() image.Rectangle
	ColorModel() color.Model
}

// pixmapRaw is used if the image's data is stored raw, as in each pixel
// stores it's own color.
type pixmapRaw struct {
	header *FileHeader
	pixels []Color
}

// pixmapMapped is used if the image's data is colormapped, e.g. 
// each pixel is a "pointer" to a value of the corresponding color map.
type pixmapMapped struct {
	header *FileHeader
	colors *ColorMap
	pixels []uint8
}

func ReadPixmap(r io.Reader, h *FileHeader, colors *ColorMap) (Pixmap, error) {
	if h.IsMapped() {
		debugf("this is a mapped pixmap")
		return readPixmapMapped(r, h, colors)
	} else {
		debugf("this is a raw pixmap")
		return readPixmapRaw(r, h)
	}
}

func readPixmapRaw(r io.Reader, h *FileHeader) (*pixmapRaw, error) {
	if h.BitsPerPixel != 24 {
		return nil, errors.New("only 24 bpp supporetd")
	}

	buf := make([]byte, 4)
	discard := make([]byte, h.BytesPerLine - h.WindowWidth * h.BitsPerPixel / 8)
	var cu uint32
	var cl Color

	pixmap := pixmapRaw{h, make([]Color, h.PixmapWidth * h.PixmapHeight)}

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
			cl = Color{
				Pixel: i, Flags: 7, Padding: 0,
				Red: uint16(((cu & h.RedMask) >> rs) << 8),
				Green: uint16(((cu & h.GreenMask) >> gs) << 8),
				Blue: uint16(((cu & h.BlueMask) >> bs) << 8),
			}
			pixmap.pixels[i] = cl
			i++
		}
		_, err := r.Read(discard) // discard line ending
		if err != nil {
			return nil, err
		}
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

func readPixmapMapped(r io.Reader, h *FileHeader, colors *ColorMap) (*pixmapMapped, error) {
	pix := pixmapMapped{}
	pix.header = h
	pix.colors = colors
	pix.pixels = make([]uint8, h.PixmapWidth * h.PixmapHeight)

	debugf("Going to read %d bytes", h.PixmapWidth * h.PixmapHeight * colormapKeySize)
	buf := make([]byte, h.PixmapWidth * h.PixmapHeight * colormapKeySize)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	var i uint32 = 0
	var x, y uint32
	for y = 0; y < h.PixmapHeight; y++ {
		for x = 0; x < h.PixmapWidth; x++ {
			pix.pixels[i] = uint8(buf[i * colormapKeySize])
			i++
		}
	}

	return &pix, nil
}

func (p *pixmapRaw) At(x, y int) color.Color {
	return &p.pixels[y * int(p.header.PixmapWidth) + x]
}

func (p *pixmapRaw) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(p.header.PixmapWidth), int(p.header.PixmapHeight))
}

func (p *pixmapRaw) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *pixmapMapped) At(x, y int) color.Color {
	id := p.pixels[y * int(p.header.PixmapWidth) + x]
	c := p.colors.Get(int(id))
	return &c
}

func (p *pixmapMapped) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(p.header.PixmapWidth), int(p.header.PixmapHeight))
}

func (p *pixmapMapped) ColorModel() color.Model {
	return color.RGBAModel
}
