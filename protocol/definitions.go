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
	NEW_LOBBY   uint8 = 0x01
	NEW_LOBBY_S uint8 = 0x01
	JOIN_LOBBY  uint8 = 0x02

	CONNECT     uint8 = 0x10
	DISCONNECT  uint8 = 0x11
	SWITCH_ROOM uint8 = 0x12

	SPRITE   uint8 = 0x20
	POSITION uint8 = 0x21
	SPEED    uint8 = 0x22
)

// 0x01 C2S
type NewLobby struct {
	GameHash []byte
}

// 0x01 S2C
type NewLobbyS struct {
	LobbyCode []byte
}

// 0x02
type JoinLobby struct {
	GameHash  []byte
	LobbyCode []byte
}

// 0x10
type Connect struct {
	Id uint16
}

// 0x11
type Disconnect struct {
	Id uint16
}

// 0x12
type SwitchRoom struct {
	Id uint16
}

// 0x20
type Sprite struct {
	Id    uint16
	Name  []byte
	Index uint8
}

// 0x21
type Position struct {
	Id        uint16
	X         uint16
	Y         uint16
	Direction uint8
}

// 0x22
type Speed struct {
	Id    uint16
	Speed uint8
}
