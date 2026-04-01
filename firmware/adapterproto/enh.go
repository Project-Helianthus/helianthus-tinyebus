package adapterproto

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedENHFrame = errors.New("enh malformed frame")
	ErrShortENHFrame     = errors.New("enh short frame")
)

type ENHCommand byte

const (
	ENHReqInit  ENHCommand = 0x0
	ENHReqSend  ENHCommand = 0x1
	ENHReqStart ENHCommand = 0x2
	ENHReqInfo  ENHCommand = 0x3

	ENHResResetted  ENHCommand = 0x0
	ENHResReceived  ENHCommand = 0x1
	ENHResStarted   ENHCommand = 0x2
	ENHResInfo      ENHCommand = 0x3
	ENHResFailed    ENHCommand = 0xA
	ENHResErrorEBUS ENHCommand = 0xB
	ENHResErrorHost ENHCommand = 0xC
)

type ENHDecodedFrame struct {
	Command     ENHCommand `json:"command"`
	CommandName string     `json:"command_name"`
	Data        byte       `json:"data"`
	Consumed    int        `json:"consumed"`
}

func (command ENHCommand) String() string {
	switch command {
	case ENHReqInit:
		return "init"
	case ENHReqSend:
		return "send"
	case ENHReqStart:
		return "start"
	case ENHReqInfo:
		return "info"
	case ENHResFailed:
		return "failed"
	case ENHResErrorEBUS:
		return "error_ebus"
	case ENHResErrorHost:
		return "error_host"
	default:
		return fmt.Sprintf("unknown(0x%02x)", byte(command))
	}
}

func EncodeENH(command ENHCommand, data byte) [2]byte {
	byte1 := byte(0xC0) | (byte(command) << 2) | ((data & 0xC0) >> 6)
	byte2 := byte(0x80) | (data & 0x3F)
	return [2]byte{byte1, byte2}
}

func EncodeENHStream(command ENHCommand, data byte) []byte {
	if command == ENHResReceived && data < 0x80 {
		return []byte{data}
	}
	encoded := EncodeENH(command, data)
	return []byte{encoded[0], encoded[1]}
}

func DecodeENH(byte1, byte2 byte) (ENHCommand, byte, error) {
	if byte1&0xC0 != 0xC0 {
		return 0, 0, fmt.Errorf("%w: invalid first byte 0x%02x", ErrMalformedENHFrame, byte1)
	}
	if byte2&0xC0 != 0x80 {
		return 0, 0, fmt.Errorf("%w: invalid second byte 0x%02x", ErrMalformedENHFrame, byte2)
	}

	command := ENHCommand((byte1 >> 2) & 0x0F)
	data := byte(((byte1 & 0x03) << 6) | (byte2 & 0x3F))
	return command, data, nil
}

func DecodeENHStream(data []byte) (ENHDecodedFrame, error) {
	if len(data) == 0 {
		return ENHDecodedFrame{}, ErrShortENHFrame
	}

	if data[0]&0x80 == 0 {
		return ENHDecodedFrame{
			Command:     ENHResReceived,
			CommandName: "received",
			Data:        data[0],
			Consumed:    1,
		}, nil
	}

	if len(data) < 2 {
		return ENHDecodedFrame{}, ErrShortENHFrame
	}

	command, payload, err := DecodeENH(data[0], data[1])
	if err != nil {
		return ENHDecodedFrame{}, err
	}

	return ENHDecodedFrame{
		Command:     command,
		CommandName: command.String(),
		Data:        payload,
		Consumed:    2,
	}, nil
}
