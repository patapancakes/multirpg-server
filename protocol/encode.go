package protocol

import (
	"encoding/binary"
	"fmt"
)

func encode(data interface{}) ([]byte, error) {
	switch data := data.(type) {
	case Connect:
		return encodeConnect(data)
	case Disconnect:
		return encodeDisconnect(data)
	case SwitchRoom:
		return encodeSwitchRoom(data)
	case Sprite:
		return encodeSprite(data)
	case Move:
		return encodeMove(data)
	case Speed:
		return encodeSpeed(data)
	default:
		return nil, fmt.Errorf("unknown message type: %T", data)
	}
}

func encodeConnect(data Connect) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{CONNECT}, id...), nil
}

func encodeDisconnect(data Disconnect) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{DISCONNECT}, id...), nil
}

func encodeSwitchRoom(data SwitchRoom) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return append([]byte{SWITCH_ROOM}, id...), nil
}

func encodeSprite(data Sprite) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	name := make([]byte, len(data.Name))
	copy(name, data.Name)

	index := make([]byte, 1)
	index[0] = byte(data.Index)

	return append(append(append(append([]byte{SPRITE}, id...), byte(len(data.Name))), name...), index...), nil
}

func encodeMove(data Move) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	x := make([]byte, 2)
	binary.LittleEndian.PutUint16(x, data.X)

	y := make([]byte, 2)
	binary.LittleEndian.PutUint16(y, data.Y)

	direction := make([]byte, 1)
	direction[0] = byte(data.Direction)

	return append(append(append(append([]byte{MOVE}, id...), x...), y...), direction...), nil
}

func encodeSpeed(data Speed) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	speed := make([]byte, 1)
	speed[0] = byte(data.Speed)

	return append(append([]byte{SPEED}, id...), speed...), nil
}
