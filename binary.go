package xwd

type Order uint8

const (
	LittleEndian Order = 0
	MSBFirst Order = 0
	BigEndian Order = 1
	LSBFirst Order = 1
	Invalid Order = 255
)

func OrderFromUint32(i uint32) Order {
	if i == 0 {
		return MSBFirst
	} else if i == 1 {
		return LSBFirst
	} else {
		return Invalid
	}
}
