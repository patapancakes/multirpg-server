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

func Decode(data []byte) (interface{}, error) {
	switch data[0] {
	case SWITCH_ROOM:
		return decodeSwitchRoom(data)
	case SPRITE:
		return decodeSprite(data)
	case POSITION:
		return decodePosition(data)
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

	return SwitchRoom{
		Id: binary.LittleEndian.Uint16(data[1:]),
	}, nil
}

func decodeSprite(data []byte) (interface{}, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	nameLength := int(data[1:2][0])
	if len(data[2:]) != nameLength+1 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Sprite{
		Name:  data[2 : 2+nameLength],
		Index: uint8(data[2+nameLength:][0]),
	}, nil
}

func decodePosition(data []byte) (interface{}, error) {
	if len(data) != 6 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Position{
		X:         binary.LittleEndian.Uint16(data[1:3]),
		Y:         binary.LittleEndian.Uint16(data[3:5]),
		Direction: uint8(data[5:][0]),
	}, nil
}

func decodeSpeed(data []byte) (interface{}, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf("invalid packet length: %d", len(data))
	}

	return Speed{
		Speed: uint8(data[1:][0]),
	}, nil
}
