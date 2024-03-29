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
	"fmt"
	"net"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Client struct {
	conn net.Conn

	terminate chan bool

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
	defer func() { c.terminate <- true }()
	for {
		select {
		case <-c.terminate:
			return
		default:
			buf := make([]byte, 300)

			n, err := c.conn.Read(buf)
			if err != nil {
				return
			}

			c.receive <- buf[:n]
		}
	}
}

func (c *Client) packetReader() {
	defer func() { c.terminate <- true }()
	for {
		select {
		case <-c.terminate:
			return
		case data, ok := <-c.receive:
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
}

func (c *Client) packetWriter() {
	defer func() { c.terminate <- true }()
	for {
		select {
		case <-c.terminate:
			return
		case data, ok := <-c.send:
			if !ok {
				return
			}

			c.conn.Write(data)
		}
	}
}

func (c *Client) sendPacket(data []byte) {
	select {
	case c.send <- data:
	default:
		c.terminate <- true
	}
}

func (c *Client) joinLobby(lobbyCode string) {
	lobby, _ := c.server.lobbies.Load(lobbyCode)

	c.lobby = lobby.(*Lobby)

	c.id = c.lobby.getFreeId()
	c.lobby.clientIds.Store(c.id, nil)
}

func (c *Client) leaveLobby() {
	c.lobby.clientIds.Delete(c.id)

	c.lobby.removeIfEmpty()

	c.lobby = nil
}

func (c *Client) joinRoom(roomId uint16) {
	room, _ := c.lobby.rooms.Load(roomId)

	c.room = room.(*Room)
	c.room.clients.Store(c, nil)

	packet, _ := protocol.Encode(protocol.ClientJoin{
		Id: c.id,
	})
	c.room.broadcast(packet, c)

	c.getRoomData()
}

func (c *Client) leaveRoom() {
	c.room.clients.Delete(c)

	packet, _ := protocol.Encode(protocol.ClientLeave{
		Id: c.id,
	})
	c.room.broadcast(packet, c)

	c.room.removeIfEmpty()

	c.room = nil
}

func (c *Client) getRoomData() {
	c.room.clients.Range(func(k, _ any) bool {
		client := k.(*Client)

		if client == c {
			return true
		}

		// Client Join
		packet, _ := protocol.Encode(protocol.ClientJoin{
			Id: client.id,
		})
		c.sendPacket(packet)

		// Sprite
		packet, _ = protocol.Encode(protocol.Sprite{
			Id:    client.id,
			Name:  client.sprite,
			Index: client.spriteIndex,
		})
		c.sendPacket(packet)

		// Position
		packet, _ = protocol.Encode(protocol.Position{
			Id:        client.id,
			X:         client.x,
			Y:         client.y,
			Direction: client.direction,
		})
		c.sendPacket(packet)

		// Speed
		packet, _ = protocol.Encode(protocol.Speed{
			Id:    client.id,
			Speed: client.speed,
		})
		c.sendPacket(packet)

		return true
	})
}

func (c *Client) disconnect() {
	if err := c.conn.Close(); err != nil {
		fmt.Printf("Connection from %s failed to close: %s\n", c.conn.RemoteAddr().String(), err)
	} else {
		fmt.Printf("Connection from %s closed\n", c.conn.RemoteAddr().String())
	}

	if c.room != nil {
		c.leaveRoom()
	}

	if c.lobby != nil {
		c.leaveLobby()
	}
}
