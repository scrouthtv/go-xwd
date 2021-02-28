package xwd

import (
	"encoding/binary"
	"image"
	"image/color"
	"image/color/palette"
	"io"
	"log"
	"errors"
	"strconv"
)

// XWDFileHeader contains information
// about an xwd image.
type XWDFileHeader struct {
	HeaderSize        uint32
	FileVersion       uint32
	PixmapFormat      uint32
	PixmapDepth       uint32
	PixmapWidth       uint32
	PixmapHeight      uint32
	XOffset           uint32
	ByteOrder         uint32
	BitmapUnit        uint32
	BitmapBitOrder    uint32
	BitmapPad         uint32
	BitsPerPixel      uint32
	BytesPerLine      uint32
	VisualClass       uint32
	RedMask           uint32
	GreenMask         uint32
	BlueMask          uint32
	BitsPerRgb        uint32
	NumberOfColors    uint32
	ColorMapEntries   uint32
	WindowWidth       uint32
	WindowHeight      uint32
	WindowX           uint32
	WindowY           uint32
	WindowBorderWidth uint32
}

// XWDColorMap is the color map of this xwd image.
// It is subject to removal.
type XWDColorMap struct {
	EntryNumber uint32
	Red         uint16
	Green       uint16
	Blue        uint16
	Flags       uint8
	Padding     uint8
}

// XWDImage groups together an xwd header and
// a paletted image.
// It's pointer implements all image functionality.
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
	xwd := XWDImage{}

	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		log.Println("no header size information")
		return nil, err
	}

	xwd.header.HeaderSize = binary.BigEndian.Uint32(buf[0:4])
	log.Println("header size: ", xwd.header.HeaderSize)

	buf = make([]byte, xwd.header.HeaderSize - 4) // read the rest of the header, we already have the first 4 bytes
	_, err = r.Read(buf)
	if err != nil {
		log.Println("short header")
		return nil, err
	}

	xwd.header.FileVersion = binary.BigEndian.Uint32(buf[4:8])
	if xwd.header.FileVersion != 7 {
		return nil, errors.New("Unsupported xwd version " + strconv.FormatUint(uint64(xwd.header.FileVersion), 10))
	}

		/*PixmapFormat:      binary.BigEndian.Uint32(buf[8:12]),
		PixmapDepth:       binary.BigEndian.Uint32(buf[12:16]),
		PixmapWidth:       binary.BigEndian.Uint32(buf[16:20]),
		PixmapHeight:      binary.BigEndian.Uint32(buf[20:24]),
		XOffset:           binary.BigEndian.Uint32(buf[24:28]),
		ByteOrder:         binary.BigEndian.Uint32(buf[28:32]),
		BitmapUnit:        binary.BigEndian.Uint32(buf[32:36]),
		BitmapBitOrder:    binary.BigEndian.Uint32(buf[36:40]),
		BitmapPad:         binary.BigEndian.Uint32(buf[40:44]),
		BitsPerPixel:      binary.BigEndian.Uint32(buf[44:48]),
		BytesPerLine:      binary.BigEndian.Uint32(buf[48:52]),
		VisualClass:       binary.BigEndian.Uint32(buf[52:56]),
		RedMask:           binary.BigEndian.Uint32(buf[56:60]),
		GreenMask:         binary.BigEndian.Uint32(buf[60:64]),
		BlueMask:          binary.BigEndian.Uint32(buf[64:68]),
		BitsPerRgb:        binary.BigEndian.Uint32(buf[68:72]),
		NumberOfColors:    binary.BigEndian.Uint32(buf[72:76]),
		ColorMapEntries:   binary.BigEndian.Uint32(buf[76:80]),
		WindowWidth:       binary.BigEndian.Uint32(buf[80:84]),
		WindowHeight:      binary.BigEndian.Uint32(buf[84:88]),
		WindowX:           binary.BigEndian.Uint32(buf[88:92]),
		WindowY:           binary.BigEndian.Uint32(buf[92:96]),
		WindowBorderWidth: binary.BigEndian.Uint32(buf[96:100]),
	}*/

	log.Println("header size:", xwd.header.HeaderSize)
	log.Println("header:", xwd.header)

	// window name
	windowName := make([]byte, xwd.header.HeaderSize-100)
	_, err = r.Read(windowName)
	if err != nil {
		log.Println("missing window size")
		return nil, err
	}

	// not used?
	colorMaps := make([]XWDColorMap, xwd.header.ColorMapEntries)
	for i := 0; i < int(xwd.header.ColorMapEntries); i++ {
		buf := make([]byte, 12)
		_, err = r.Read(buf)
		if err != nil {
			return nil, err
		}
		colorMaps[i] = XWDColorMap{
			EntryNumber: binary.BigEndian.Uint32(buf[0:4]),
			Red:         binary.BigEndian.Uint16(buf[4:6]),
			Green:       binary.BigEndian.Uint16(buf[6:8]),
			Blue:        binary.BigEndian.Uint16(buf[8:10]),
			Flags:       uint8(buf[10]),
			Padding:     uint8(buf[11]),
		}
	}

	rect := image.Rect(0, 0, int(xwd.header.PixmapWidth), int(xwd.header.PixmapHeight))

	// create PalettedImage
	xwd.image = image.NewPaletted(rect, palette.WebSafe)
	for x := 0; x < int(xwd.header.PixmapHeight); x++ {
		for y := 0; y < int(xwd.header.PixmapWidth); y++ {
			buf := make([]byte, 4)
			_, err = r.Read(buf)
			if err != nil {
				log.Println("error reading pixel @", x, y)
				return nil, err
			}
			if x == 5 && y == 5 {
				log.Println("At position 5/5:", uint8(buf[2]), uint8(buf[1]), uint8(buf[0]))
			}
			xwd.image.Set(y, x, color.RGBA{
				R: uint8(buf[2]),
				G: uint8(buf[1]),
				B: uint8(buf[0]),
				A: 255,
			})
		}
	}
	return &xwd, nil
}
