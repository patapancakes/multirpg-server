package protocol

import (
	"encoding/binary"
	"fmt"
)

func decode(data []byte) (interface{}, error) {
	switch data[0] {
	case SWITCH_ROOM:
		return decodeSwitchRoom(data)
	case SPRITE:
		return decodeSprite(data)
	case MOVE:
		return decodeMove(data)
	case SPEED:
		return decodeSpeed(data)
	default:
		return nil, fmt.Errorf("unknown packet type: %d", data[0])
	}
}

func decodeSwitchRoom(data []byte) (interface{}, error) {
	if len(data) != 3 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return SwitchRoom{Id: binary.LittleEndian.Uint16(data[1:])}, nil
}

func decodeSprite(data []byte) (interface{}, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	nameLength := int(data[3:4][0])
	if len(data[4:]) != nameLength + 1 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Sprite{Id: binary.LittleEndian.Uint16(data[1:3]), Name: data[4:4+nameLength], Index: uint8(data[4+nameLength:][0])}, nil
}

func decodeMove(data []byte) (interface{}, error) {
	if len(data) != 8 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Move{Id: binary.LittleEndian.Uint16(data[1:3]), X: binary.LittleEndian.Uint16(data[3:5]), Y: binary.LittleEndian.Uint16(data[5:7]), Direction: uint8(data[7:][0])}, nil
}

func decodeSpeed(data []byte) (interface{}, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Speed{Id: binary.LittleEndian.Uint16(data[1:3]), Speed: uint8(data[3:][0])}, nil
}
