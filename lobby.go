package main

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

type Lobby struct {
	gameHash []byte

	rooms     map[uint16]*Room
	clientIds map[uint16]bool
}

func (s *Server) createLobby(gameHash []byte) *Lobby {
	lobby := &Lobby{
		gameHash:  gameHash,

		rooms:     make(map[uint16]*Room),
		clientIds: make(map[uint16]bool),
	}

	// Room 0 is the room clients are put in when they first join the lobby
	// Clients are expected to send a switch room packet to join a game room
	lobby.rooms[0] = lobby.createRoom(0)

	return lobby
}

func (l *Lobby) getFreeId() uint16 {
	for i := uint16(0); i < 0xFFFF; i++ {
		if _, ok := l.clientIds[i]; !ok {
			return i
		}
	}

	// This should never happen, if it does then somehow all ids are being used
	return 0
}
