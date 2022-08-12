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

import (
	"net"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Client struct {
	conn net.Conn
	id   uint16

	room *Room

	sprite      []byte
	spriteIndex uint8

	x         uint16
	y         uint16
	direction uint8

	speed uint8
}

// Listen for incoming packets from the client
func (c *Client) listen() {
	for {
		buf := make([]byte, 300)

		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		packet := &Packet{
			sender: c,
			data:   buf[:n],
		}

		go packet.process()
	}
}

func (c *Client) getRoomData() {
	for client := range c.room.clients {
		if client == c {
			continue
		}

		// Connect
		packet, _ := protocol.Encode(protocol.Connect{
			Id: client.id,
		})
		c.conn.Write(packet)

		// Sprite
		packet, _ = protocol.Encode(protocol.Sprite{
			Id:    client.id,
			Name:  client.sprite,
			Index: client.spriteIndex,
		})
		c.conn.Write(packet)

		// Position
		packet, _ = protocol.Encode(protocol.Position{
			Id:        client.id,
			X:         client.x,
			Y:         client.y,
			Direction: client.direction,
		})
		c.conn.Write(packet)
	}
}

func (c *Client) joinRoom() {
	// Redundant almost always, do not send room data or connect packets if in the default room
	if c.room.id == 0 {
		return
	}

	c.getRoomData()

	packet, _ := protocol.Encode(protocol.Connect{
		Id: c.id,
	})
	c.room.broadcast(packet, c)
}

func (c *Client) leaveRoom() {
	delete(c.room.server.rooms[c.room.id].clients, c)

	packet, _ := protocol.Encode(protocol.Disconnect{
		Id: c.id,
	})
	c.room.broadcast(packet, c)
}
