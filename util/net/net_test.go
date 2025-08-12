package net

import "testing"

func TestGetLocalIPV4(t *testing.T) {
	ip, err := GetLocalIPV4()
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}
