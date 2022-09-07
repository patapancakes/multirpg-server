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

func Encode(packet any) ([]byte, error) {
	switch packet := packet.(type) {
	case Connect:
		return packSegments(2, CONNECT, packet.Id), nil
	case Disconnect:
		return packSegments(2, DISCONNECT, packet.Id), nil
	case SwitchRoom:
		return packSegments(2, SWITCH_ROOM, packet.Id), nil
	case Sprite:
		return packSegments(0, SPRITE, packet.Id, uint8(len(packet.Name)), packet.Name, packet.Index), nil
	case Position:
		return packSegments(5, POSITION, packet.Id, packet.X, packet.Y, packet.Direction), nil
	case Speed:
		return packSegments(3, SPEED, packet.Id, packet.Speed), nil
	default:
		return nil, fmt.Errorf("unknown packet type: %T", packet)
	}
}

func packSegments(length int, segments ...any) []byte {
	buf := make([]byte, length)
	for _, segment := range segments {
		switch segment := segment.(type) {
		case byte:
			buf = append(buf, segment)
		case []byte:
			buf = append(buf, segment...)
		case uint16:
			ibuf := make([]byte, 2)
			binary.LittleEndian.PutUint16(ibuf, segment)

			buf = append(buf, ibuf...)
		}
	}

	return buf
}
