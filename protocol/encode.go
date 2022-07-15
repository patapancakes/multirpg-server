package protocol

import (
	"encoding/binary"
	"fmt"
)

func Encode(data interface{}) ([]byte, error) {
	switch data := data.(type) {
	case Connect:
		return encodeConnect(data)
	case Disconnect:
		return encodeDisconnect(data)
	case SwitchRoom:
		return encodeSwitchRoom(data)
	case Sprite:
		return encodeSprite(data)
	case Position:
		return encodePosition(data)
	case Speed:
		return encodeSpeed(data)
	default:
		return nil, fmt.Errorf("unknown packet type: %T", data)
	}
}

func encodeConnect(data Connect) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{
		CONNECT},
		id...
	), nil
}

func encodeDisconnect(data Disconnect) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{
		DISCONNECT},
		id...
	), nil
}

func encodeSwitchRoom(data SwitchRoom) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{
		SWITCH_ROOM},
		id...
	), nil
}

func encodeSprite(data Sprite) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append(append(append(append([]byte{
		SPRITE},
		id...),
		uint8(len(data.Name))),
		data.Name...),
		data.Index,
	), nil
}

func encodePosition(data Position) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	x := make([]byte, 2)
	binary.LittleEndian.PutUint16(x, data.X)

	y := make([]byte, 2)
	binary.LittleEndian.PutUint16(y, data.Y)

	return append(append(append(append([]byte{
		POSITION},
		id...),
		x...),
		y...),
		data.Direction,
	), nil
}

func encodeSpeed(data Speed) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append(append([]byte{
		SPEED},
		id...),
		data.Speed,
	), nil
}
