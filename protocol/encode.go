package protocol

/*
multirpg-server
https://github.com/Gamizard/multirpg-server

Copyright (C) 2022 azarashi <azarashi@majestaria.fun>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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

	return combine(CONNECT, id), nil
}

func encodeDisconnect(data Disconnect) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return combine(DISCONNECT, id), nil
}

func encodeSwitchRoom(data SwitchRoom) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return combine(SWITCH_ROOM, id), nil
}

func encodeSprite(data Sprite) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return combine(SPRITE, id, uint8(len(data.Name)), data.Name, data.Index), nil
}

func encodePosition(data Position) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	x := make([]byte, 2)
	binary.LittleEndian.PutUint16(x, data.X)

	y := make([]byte, 2)
	binary.LittleEndian.PutUint16(y, data.Y)

	return combine(POSITION, id, x, y, data.Direction), nil
}

func encodeSpeed(data Speed) ([]byte, error) {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, data.Id)

	return combine(SPEED, id, data.Speed), nil
}

// combine combines serveral bytes or byte arrays into a single byte array
func combine(segments ...any) []byte {
	var buf []byte
	for _, segment := range segments {
		switch segment := segment.(type) {
		case byte:
			buf = append(buf, segment)
		case []byte:
			buf = append(buf, segment...)
		}
	}

	return buf
}
