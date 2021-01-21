package drivehash

import (
	"fmt"
	"testing"
)

func TestNewMultiHash(t *testing.T) {
	hasher := NewMultiHash()

	hasher.Write([]byte("hello world"))

	if fmt.Sprintf("%x", hasher.MD5()) != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Fatal("MD5 mismatch")
	}

	if fmt.Sprintf("%x", hasher.SHA1()) != "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed" {
		t.Fatal("SHA1 mismatch")
	}

	if fmt.Sprintf("%x", hasher.SHA256()) != "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" {
		t.Fatal("SHA256 mismatch")
	}

	hasher.Reset()
}
