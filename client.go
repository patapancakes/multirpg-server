package main

import (
	"net"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Client struct {
	conn net.Conn
	id uint16

	room *Room

	sprite []byte
	spriteIndex uint8

	x uint16
	y uint16
	direction uint8

	speed uint8
}

func (c *Client) listen() {
	for	{
		buf := make([]byte, 300)

		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		packet := &Packet{
			sender: c,
			data: buf[:n],
		}

		go packet.process()
	}
}

func (c *Client) getRoomData() {
	for otherClient := range c.room.clients {
		if otherClient == c {
			continue
		}

		// Connect
		packet, _ := protocol.Encode(protocol.Connect{
			Id: otherClient.id,
		})
		c.conn.Write(packet)

		// Sprite
		packet, _ = protocol.Encode(protocol.Sprite{
			Id: otherClient.id,
			Name: otherClient.sprite,
			Index: otherClient.spriteIndex,
		})
		c.conn.Write(packet)

		// Position
		packet, _ = protocol.Encode(protocol.Position{
			Id: otherClient.id,
			X: otherClient.x,
			Y: otherClient.y,
			Direction: otherClient.direction,
		})
		c.conn.Write(packet)

		// Speed
		packet, _ = protocol.Encode(protocol.Speed{
			Id: otherClient.id,
			Speed: otherClient.speed,
		})
		c.conn.Write(packet)
	}
}

func (c *Client) joinRoom() {
	// Redundant almost always, do not send room data or connect packets if on the title screen
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
