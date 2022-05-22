package exdhcp

import "testing"

func TestRandomMAC(t *testing.T) {
	for i := 0; i < 5; i++ {
		hwaddr := randomMAC()
		t.Log(hwaddr.String())
	}
}
