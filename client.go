package main

import (
	"fmt"
	"net"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Client struct {
	conn net.Conn
	room *Room
	id uint16

	sprite []byte
	spriteIndex uint8

	x uint16
	y uint16
	direction uint8

	speed uint8
}

func (c *Client) listen() {
	for	{
		buf := make([]byte, 1024)

		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		packet := &Packet{
			sender: c,
			data: buf[:n],
		}

		go packet.recieve()
	}
}

func (c *Client) sendRoomData() {
	for _, otherClient := range c.room.clients {
		if otherClient.id == c.id {
			continue
		}

		// Connect
		packet, err := protocol.Encode(protocol.Connect{Id: otherClient.id})
		if err != nil {
			fmt.Println(err)
		}

		c.conn.Write(packet)

		// Sprite
		packet, err = protocol.Encode(protocol.Sprite{
			Id: otherClient.id,
			Name: otherClient.sprite,
			Index: otherClient.spriteIndex,
		})
		if err != nil {
			fmt.Println(err)
		}

		c.conn.Write(packet)

		// Position
		packet, err = protocol.Encode(protocol.Move{
			Id: otherClient.id,
			X: otherClient.x,
			Y: otherClient.y,
			Direction: otherClient.direction,
		})
		if err != nil {
			fmt.Println(err)
		}

		c.conn.Write(packet)

		// Speed
		packet, err = protocol.Encode(protocol.Speed{
			Id: otherClient.id,
			Speed: otherClient.speed,
		})
		if err != nil {
			fmt.Println(err)
		}

		c.conn.Write(packet)
	}
}

func (c *Client) handleConnect() {
	c.sendRoomData()

	packet, err := protocol.Encode(protocol.Connect{Id: c.id})
	if err != nil {
		fmt.Println(err)
	}

	c.room.broadcast(packet)
}

func (c *Client) handleDisconnect() {
	delete(c.room.server.rooms[c.room.id].clients, c.id)

	packet, err := protocol.Encode(protocol.Disconnect{Id: c.id})
	if err != nil {
		fmt.Println(err)
	}

	c.room.broadcast(packet)
}
