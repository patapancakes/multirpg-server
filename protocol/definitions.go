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

const (
	CONNECT     uint8 = 0x01
	DISCONNECT  uint8 = 0x02
	SWITCH_ROOM uint8 = 0x03
	SPRITE      uint8 = 0x10
	POSITION    uint8 = 0x11
	SPEED       uint8 = 0x12
)

//0x01
type Connect struct {
	Id uint16
}

//0x02
type Disconnect struct {
	Id uint16
}

//0x03
type SwitchRoom struct {
	Id uint16
}

//0x10
type Sprite struct {
	Id    uint16
	Name  []byte
	Index uint8
}

//0x11
type Position struct {
	Id        uint16
	X         uint16
	Y         uint16
	Direction uint8
}
