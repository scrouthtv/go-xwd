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

const (
	xwdVersion = 7
	xwdHeaderSize = 100 // size of the header without the window name

	pixmapFormat = 2 // ZPixmap
	xOffset = 0 // "number of pixels offset in X direction", idk
	bitmapPad = 32 // "8, 16, 32 either XY or ZPixmap", idk
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
	ByteOrder         Order
	BitmapUnit        uint32
	BitmapBitOrder    Order
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
	WindowName string
}

type Order uint8

const (
	BigEndian Order = 0
	LSBFirst Order = 0
	LittleEndian Order = 1
	MSBFirst Order = 1
	Invalid Order = 255
)

func OrderFromUint32(i uint32) Order {
	if i == 0 {
		return LSBFirst
	} else if i == 1 {
		return MSBFirst
	} else {
		return Invalid
	}
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
	xwd := XWDImage{}

	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		log.Println("no header size information")
		return nil, err
	}

	xwd.header.HeaderSize = binary.BigEndian.Uint32(buf[0:4])
	log.Println("header size: ", xwd.header.HeaderSize)

	buf = make([]byte, xwdHeaderSize - 4) // read the rest of the header, we already have the first value
	_, err = r.Read(buf)
	if err != nil {
		log.Println("short header")
		return nil, err
	}

	xwd.header.FileVersion = binary.BigEndian.Uint32(buf[0:4])
	if xwd.header.FileVersion != xwdVersion {
		return nil, errors.New("Unsupported xwd version " + strconv.FormatUint(uint64(xwd.header.FileVersion), 10))
	}

	xwd.header.PixmapFormat = binary.BigEndian.Uint32(buf[4:8])
	if xwd.header.FileVersion != pixmapFormat {
		return nil, errors.New("Unsupported pixmap format " + strconv.FormatUint(uint64(xwd.header.PixmapFormat), 10))
	}

	xwd.header.PixmapDepth = binary.BigEndian.Uint32(buf[8:12])
	xwd.header.PixmapWidth = binary.BigEndian.Uint32(buf[12:16])
	xwd.header.PixmapHeight = binary.BigEndian.Uint32(buf[16:20])

	xwd.header.XOffset = binary.BigEndian.Uint32(buf[20:24])
	if xwd.header.XOffset != xOffset {
		return nil, errors.New("Unsupported xoffset " + strconv.FormatUint(uint64(xwd.header.XOffset), 10))
	}

	xwd.header.ByteOrder = OrderFromUint32(binary.BigEndian.Uint32(buf[24:28]))
	if xwd.header.ByteOrder == Invalid {
		return nil, errors.New("Unsupported byte order " + strconv.FormatUint(uint64(binary.BigEndian.Uint32(buf[24:28])), 10))
	}

	xwd.header.BitmapUnit = binary.BigEndian.Uint32(buf[28:32])

	xwd.header.BitmapBitOrder = OrderFromUint32(binary.BigEndian.Uint32(buf[32:36]))
	if xwd.header.BitmapBitOrder == Invalid {
		return nil, errors.New("Unsupported bit order " + strconv.FormatUint(uint64(binary.BigEndian.Uint32(buf[32:36])), 10))
	}

	xwd.header.BitmapPad = binary.BigEndian.Uint32(buf[36:40])
	if xwd.header.BitmapPad != bitmapPad {
		return nil, errors.New("Unsupported bitmapd pad " + strconv.FormatUint(uint64(xwd.header.BitmapPad), 10))
	}

	xwd.header.BitsPerPixel = binary.BigEndian.Uint32(buf[40:44])
	xwd.header.BytesPerLine = binary.BigEndian.Uint32(buf[44:48])
	xwd.header.VisualClass = binary.BigEndian.Uint32(buf[48:52])

	xwd.header.RedMask = binary.BigEndian.Uint32(buf[52:56])
	if xwd.header.RedMask == 0 {
		return nil, errors.New("Red mask invalid")
	}
	xwd.header.GreenMask = binary.BigEndian.Uint32(buf[56:60])
	if xwd.header.GreenMask == 0 {
		return nil, errors.New("Green mask invalid")
	}
	xwd.header.BlueMask = binary.BigEndian.Uint32(buf[60:64])
	if xwd.header.BlueMask == 0 {
		return nil, errors.New("Blue mask invalid")
	}

	xwd.header.BitsPerRgb = binary.BigEndian.Uint32(buf[64:68])
	xwd.header.NumberOfColors = binary.BigEndian.Uint32(buf[68:72])
	xwd.header.ColorMapEntries = binary.BigEndian.Uint32(buf[72:76])
	xwd.header.WindowWidth = binary.BigEndian.Uint32(buf[76:80])
	xwd.header.WindowHeight = binary.BigEndian.Uint32(buf[80:84])
	xwd.header.WindowX = binary.BigEndian.Uint32(buf[84:88])
	xwd.header.WindowY = binary.BigEndian.Uint32(buf[88:92])
	xwd.header.WindowBorderWidth = binary.BigEndian.Uint32(buf[92:96])

	log.Println("header:", xwd.header)

	// window name
	buf = make([]byte, xwd.header.HeaderSize - xwdHeaderSize)
	_, err = r.Read(buf)
	if err != nil {
		log.Println("error reading window name")
		return nil, err
	}

	// strip the null terminator:
	end := xwd.header.HeaderSize - xwdHeaderSize - 1
	if end < 0 {
		end = 0
	}
	xwd.header.WindowName = string(buf[:end])

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
