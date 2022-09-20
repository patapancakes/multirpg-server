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

import "sync"

type Room struct {
	lobby *Lobby
	id    uint16

	clients sync.Map
}

func (l *Lobby) createRoom(id uint16) *Room {
	return &Room{
		lobby:   l,
		id:      id,
	}
}

func (r *Room) broadcast(data []byte, sender *Client) {
	r.clients.Range(func(k, _ any) bool {
		client := k.(*Client)

		if client == sender {
			return true
		}

		client.sendPacket(data)

		return true
	})
}

func (r *Room) removeIfEmpty() {
	var hasPlayers bool

	r.clients.Range(func(_, _ any) bool {
		hasPlayers = true

		return false
	})

	if hasPlayers {
		return
	}

	r.lobby.rooms.Delete(r.id)
}
