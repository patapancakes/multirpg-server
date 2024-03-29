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

package main

import (
	"math/rand"
	"sync"
)

type Lobby struct {
	server *Server

	gameHash []byte

	rooms     sync.Map
	clientIds sync.Map
}

func generateLobbyCode() string {
	const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const runesLen = len(runes)

	b := make([]byte, 6)
	for i := range b {
		b[i] = runes[rand.Intn(runesLen)]
	}

	return string(b)
}

func (s *Server) createLobby(gameHash []byte) *Lobby {
	return &Lobby{
		server:   s,
		gameHash: gameHash,
	}
}

func (l *Lobby) getFreeId() uint16 {
	for i := uint16(0); i < 0xFFFF; i++ {
		if _, ok := l.clientIds.Load(i); !ok {
			return i
		}
	}

	// This should never happen, if it does then somehow all ids are being used
	return 0
}

func (l *Lobby) removeIfEmpty() {
	var hasRooms bool

	l.rooms.Range(func(_, _ any) bool {
		hasRooms = true

		return false
	})

	if hasRooms {
		return
	}

	l.server.lobbies.Delete(l.gameHash)
}
