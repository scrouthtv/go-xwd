package xwd

import "testing"

func TestShift(t *testing.T) {
	if shiftwidth(255) != 0 {
		t.Errorf("shiftwidth(255) = %d != 0", shiftwidth(255))
	}
	if shiftwidth(65280) != 8 {
		t.Errorf("shiftwidth(65280) = %d != 8", shiftwidth(62580))
	}
	if shiftwidth(16711680) != 16 {
		t.Errorf("shiftwidth(16711680) = %d != 16", shiftwidth(16711680))
	}
}
