package xwd

import "fmt"
import "strconv"

// UnsupportedError is returned if a key is not supported.
type UnsupportedError struct {
	key string
	value string
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("Invalid %s: %s", e.key, e.value)
}

func i32toa(in uint32) string {
	return strconv.FormatUint(uint64(in), 10)
}
