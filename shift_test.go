package xwd

import "testing"

func TestShift(t *testing.T) {
	t.Log(shiftwidth(255))
	t.Log(shiftwidth(65280))
	t.Log(shiftwidth(16711680))
}
