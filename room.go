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

type Room struct {
	server *Server
	id     uint16

	clients map[*Client]bool
}

func (s *Server) createRoom(id uint16) *Room {
	return &Room{
		server:  s,
		id:      id,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) broadcast(data []byte, sender *Client) {
	// Packet broadcasting is disabled while on the title screen
	if r.id == 0 {
		return
	}

	for client := range r.clients {
		if client == sender {
			continue
		}

		client.conn.Write(data)
	}
}
