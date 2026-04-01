package adapterproto

import "testing"

func TestDecodeENHStreamShortFormReceived(t *testing.T) {
	t.Parallel()

	frame, err := DecodeENHStream([]byte{0x5A})
	if err != nil {
		t.Fatalf("DecodeENHStream() error = %v", err)
	}
	if frame.Command != ENHResReceived {
		t.Fatalf("Command = %v; want %v", frame.Command, ENHResReceived)
	}
	if frame.Data != 0x5A {
		t.Fatalf("Data = 0x%02x; want 0x5a", frame.Data)
	}
	if frame.Consumed != 1 {
		t.Fatalf("Consumed = %d; want 1", frame.Consumed)
	}
}

func TestEncodeDecodeENHStreamRoundTrip(t *testing.T) {
	t.Parallel()

	encoded := EncodeENHStream(ENHReqInfo, 0x06)
	if len(encoded) != 2 {
		t.Fatalf("EncodeENHStream() len = %d; want 2", len(encoded))
	}
	frame, err := DecodeENHStream(encoded)
	if err != nil {
		t.Fatalf("DecodeENHStream() error = %v", err)
	}
	if frame.Command != ENHReqInfo {
		t.Fatalf("Command = %v; want %v", frame.Command, ENHReqInfo)
	}
	if frame.Data != 0x06 {
		t.Fatalf("Data = 0x%02x; want 0x06", frame.Data)
	}
	if frame.Consumed != 2 {
		t.Fatalf("Consumed = %d; want 2", frame.Consumed)
	}

	reencoded := EncodeENHStream(frame.Command, frame.Data)
	if len(reencoded) != len(encoded) {
		t.Fatalf("reencoded len = %d; want %d", len(reencoded), len(encoded))
	}
	for idx := range encoded {
		if encoded[idx] != reencoded[idx] {
			t.Fatalf("reencoded[%d] = 0x%02x; want 0x%02x", idx, reencoded[idx], encoded[idx])
		}
	}
}

func TestDecodeENHLowLevel(t *testing.T) {
	t.Parallel()

	byte1, byte2 := EncodeENH(ENHReqSend, 0x9A)[0], EncodeENH(ENHReqSend, 0x9A)[1]
	command, data, err := DecodeENH(byte1, byte2)
	if err != nil {
		t.Fatalf("DecodeENH() error = %v", err)
	}
	if command != ENHReqSend {
		t.Fatalf("Command = %v; want %v", command, ENHReqSend)
	}
	if data != 0x9A {
		t.Fatalf("Data = 0x%02x; want 0x9a", data)
	}
}
