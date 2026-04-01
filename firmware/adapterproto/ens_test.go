package adapterproto

import (
	"bytes"
	"testing"
)

func TestEncodeDecodeENSRoundTrip(t *testing.T) {
	t.Parallel()

	input := []byte{0x11, ENSByteEscape, ENSByteSync, 0x7F}
	encoded := EncodeENS(input)
	decoded, err := DecodeENS(encoded)
	if err != nil {
		t.Fatalf("DecodeENS() error = %v", err)
	}
	if !bytes.Equal(input, decoded) {
		t.Fatalf("decoded = %x; want %x", decoded, input)
	}
}

func TestDecodeENSRejectsDanglingEscape(t *testing.T) {
	t.Parallel()

	_, err := DecodeENS([]byte{ENSByteEscape})
	if err == nil {
		t.Fatalf("DecodeENS() error = nil; want error")
	}
}
