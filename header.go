package xwd

import (
	"encoding/binary"
	"image"
	"image/color"
	"io"
)

const (
	xwdVersion    = 7
	xwdHeaderSize = 100 // size of the header without the window name

	pixmapFormat = 2 // ZPixmap
	xOffset      = 0 // "number of pixels offset in X direction", idk
)

// FileHeader contains information
// about an xwd image.
type FileHeader struct {
	WindowName        string
	HeaderSize        uint32
	FileVersion       uint32
	PixmapFormat      uint32 /* XYBitmap, XYPixmap, ZPixmap */
	PixmapDepth       uint32
	PixmapWidth       uint32
	PixmapHeight      uint32
	XOffset           uint32
	BitmapUnit        uint32
	BitmapPad         uint32 /* https://www.youtube.com/watch?v=dGt0Y1q0M1A */
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
	ByteOrder         Order
	BitmapBitOrder    Order
}

// IsMapped returns whether this image's data is colormapped
// or written directly in raw into the image.
func (h *FileHeader) IsMapped() bool {
	if h.ColorMapEntries == 0 {
		return false
	}
	return h.NumberOfColors == h.ColorMapEntries
}

func (h *FileHeader) Config() image.Config {
	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      int(h.PixmapWidth),
		Height:     int(h.PixmapHeight),
	}
}

// ReadHeader reads the header of an xwd image from r and returns the header or any error.
// The entire header gets read from the reader.
func ReadHeader(r io.Reader) (*FileHeader, error) {
	header := FileHeader{}

	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return nil, &IOError{err, "reading header size"}
	}

	header.HeaderSize = binary.BigEndian.Uint32(buf[0:4])
	debugf("header size: %d", header.HeaderSize)

	buf = make([]byte, xwdHeaderSize-4) // read the rest of the header, we already have the first value
	_, err = r.Read(buf)
	if err != nil {
		return nil, &IOError{err, "reading header"}
	}

	header.FileVersion = binary.BigEndian.Uint32(buf[0:4])
	if header.FileVersion != xwdVersion {
		return nil, &UnsupportedError{"xwd version", i32toa(header.FileVersion)}
	}

	header.PixmapFormat = binary.BigEndian.Uint32(buf[4:8])
	if header.PixmapFormat != pixmapFormat {
		return nil, &UnsupportedError{"pixmap format", i32toa(header.PixmapFormat)}
	}

	header.PixmapDepth = binary.BigEndian.Uint32(buf[8:12])
	header.PixmapWidth = binary.BigEndian.Uint32(buf[12:16])
	header.PixmapHeight = binary.BigEndian.Uint32(buf[16:20])

	header.XOffset = binary.BigEndian.Uint32(buf[20:24])
	if header.XOffset != xOffset {
		return nil, &UnsupportedError{"xoffset", i32toa(header.XOffset)}
	}

	header.ByteOrder = OrderFromUint32(binary.BigEndian.Uint32(buf[24:28]))
	if header.ByteOrder != BigEndian {
		return nil, &UnsupportedError{"byte order", i32toa(binary.BigEndian.Uint32(buf[24:28]))}
	}

	header.BitmapUnit = binary.BigEndian.Uint32(buf[28:32])

	header.BitmapBitOrder = OrderFromUint32(binary.BigEndian.Uint32(buf[32:36]))
	if header.BitmapBitOrder != BigEndian {
		return nil, &UnsupportedError{"bit order", i32toa(binary.BigEndian.Uint32(buf[32:36]))}
	}

	header.BitmapPad = binary.BigEndian.Uint32(buf[36:40])
	header.BitsPerPixel = binary.BigEndian.Uint32(buf[40:44])
	header.BytesPerLine = binary.BigEndian.Uint32(buf[44:48])
	header.VisualClass = binary.BigEndian.Uint32(buf[48:52])

	// masks are only set if the image is not colormapped
	header.RedMask = binary.BigEndian.Uint32(buf[52:56])
	header.GreenMask = binary.BigEndian.Uint32(buf[56:60])
	header.BlueMask = binary.BigEndian.Uint32(buf[60:64])

	header.BitsPerRgb = binary.BigEndian.Uint32(buf[64:68])
	header.NumberOfColors = binary.BigEndian.Uint32(buf[68:72])
	header.ColorMapEntries = binary.BigEndian.Uint32(buf[72:76])
	header.WindowWidth = binary.BigEndian.Uint32(buf[76:80])
	header.WindowHeight = binary.BigEndian.Uint32(buf[80:84])
	header.WindowX = binary.BigEndian.Uint32(buf[84:88])
	header.WindowY = binary.BigEndian.Uint32(buf[88:92])
	header.WindowBorderWidth = binary.BigEndian.Uint32(buf[92:96])

	// window name
	buf = make([]byte, header.HeaderSize-xwdHeaderSize)
	_, err = r.Read(buf)
	if err != nil {
		return nil, &IOError{err, "reading window name"}
	}

	// strip the null terminator:
	end := header.HeaderSize - xwdHeaderSize
	if end > 0 {
		end--
	}
	header.WindowName = string(buf[:end])

	return &header, nil
}
