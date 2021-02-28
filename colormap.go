package xwd

import (
)

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
