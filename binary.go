package xwd

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
