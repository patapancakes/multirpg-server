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
	"fmt"
	"net"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Client struct {
	conn net.Conn

	send    chan []byte
	receive chan []byte

	server *Server
	lobby  *Lobby
	room   *Room

	id uint16

	sprite      []byte
	spriteIndex uint8

	x         uint16
	y         uint16
	direction uint8
	speed     uint8
}

// Listen for incoming packets from the client
func (c *Client) listen() {
	defer c.closeConn()
	for {
		buf := make([]byte, 300)

		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		c.receive <- buf[:n]
	}
}

func (c *Client) packetReader() {
	defer c.closeConn()
	for {
		data, ok := <-c.receive
		if !ok {
			return
		}

		if len(data) < 1 {
			continue
		}

		packet := &Packet{
			sender: c,
			data:   data,
		}

		packet.process()
	}
}

func (c *Client) packetWriter() {
	defer c.closeConn()
	for {
		data, ok := <-c.send
		if !ok {
			return
		}

		c.conn.Write(data)
	}
}

func (c *Client) joinLobby(lobbyCode string) {
	c.lobby = c.server.lobbies[lobbyCode]

	c.id = c.lobby.getFreeId()
	c.lobby.clientIds[c.id] = true
}

func (c *Client) leaveLobby() {
	delete(c.lobby.clientIds, c.id)

	c.lobby = nil
}

func (c *Client) joinRoom(roomId uint16) {
	c.room = c.lobby.rooms[roomId]
	c.lobby.rooms[roomId].clients[c] = true

	packet, _ := protocol.Encode(protocol.ClientJoin{
		Id: c.id,
	})
	c.room.broadcast(packet, c)

	c.getRoomData()
}

func (c *Client) leaveRoom() {
	delete(c.lobby.rooms[c.room.id].clients, c)

	packet, _ := protocol.Encode(protocol.ClientLeave{
		Id: c.id,
	})
	c.room.broadcast(packet, c)

	c.room.removeIfEmpty()

	c.room = nil
}

func (c *Client) getRoomData() {
	for client := range c.room.clients {
		if client == c {
			continue
		}

		// Client Join
		packet, _ := protocol.Encode(protocol.ClientJoin{
			Id: client.id,
		})
		c.send <- packet

		// Sprite
		packet, _ = protocol.Encode(protocol.Sprite{
			Id:    client.id,
			Name:  client.sprite,
			Index: client.spriteIndex,
		})
		c.send <- packet

		// Position
		packet, _ = protocol.Encode(protocol.Position{
			Id:        client.id,
			X:         client.x,
			Y:         client.y,
			Direction: client.direction,
		})
		c.send <- packet

		// Speed
		packet, _ = protocol.Encode(protocol.Speed{
			Id:    client.id,
			Speed: client.speed,
		})
		c.send <- packet
	}
}

func (c *Client) closeConn() {
	if err := c.conn.Close(); err != nil {
		if err != net.ErrClosed {
			fmt.Printf("Connection from %s failed to close: %s\n", c.conn.RemoteAddr().String(), err)
		}
	} else {
		fmt.Printf("Connection from %s closed\n", c.conn.RemoteAddr().String())
	}
}

func (c *Client) disconnect() {
	close(c.send)
	close(c.receive)

	if c.room != nil {
		c.leaveRoom()
	}

	if c.lobby != nil {
		c.leaveLobby()
	}
}
