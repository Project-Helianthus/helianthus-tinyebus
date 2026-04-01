package adapterproto

import (
	"errors"
	"fmt"
)

var ErrMalformedENSFrame = errors.New("ens malformed frame")

type ENSCommand byte

const (
	ENSCommandData ENSCommand = 0x1
)

const (
	ENSByteEscape = byte(0xA9)
	ENSByteSync   = byte(0xAA)
)

func EncodeENS(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	encoded := make([]byte, 0, len(data))
	for _, payloadByte := range data {
		switch payloadByte {
		case ENSByteEscape:
			encoded = append(encoded, ENSByteEscape, 0x00)
		case ENSByteSync:
			encoded = append(encoded, ENSByteEscape, 0x01)
		default:
			encoded = append(encoded, payloadByte)
		}
	}
	return encoded
}

func DecodeENS(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	decoded := make([]byte, 0, len(data))
	escapePending := false
	for _, payloadByte := range data {
		if escapePending {
			escapePending = false
			switch payloadByte {
			case 0x00:
				decoded = append(decoded, ENSByteEscape)
			case 0x01:
				decoded = append(decoded, ENSByteSync)
			default:
				return nil, fmt.Errorf("%w: invalid escaped byte 0x%02x", ErrMalformedENSFrame, payloadByte)
			}
			continue
		}

		if payloadByte == ENSByteEscape {
			escapePending = true
			continue
		}
		if payloadByte == ENSByteSync {
			return nil, fmt.Errorf("%w: unexpected sync byte 0x%02x", ErrMalformedENSFrame, payloadByte)
		}
		decoded = append(decoded, payloadByte)
	}

	if escapePending {
		return nil, fmt.Errorf("%w: dangling escape byte", ErrMalformedENSFrame)
	}
	return decoded, nil
}
