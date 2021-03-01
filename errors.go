package xwd

import (
	"fmt"
	"strconv"
)

// UnsupportedError is returned if a key is not supported.
type UnsupportedError struct {
	key   string
	value string
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.key, e.value)
}

func i32toa(in uint32) string {
	return strconv.FormatUint(uint64(in), 10)
}

// IOError is returned if something could not be read / written.
type IOError struct {
	err  error
	step string
}

func (e *IOError) Error() string {
	return fmt.Sprintf("error %s: %s", e.step, e.err.Error())
}

func (e *IOError) Unwrap() error {
	return e.err
}
