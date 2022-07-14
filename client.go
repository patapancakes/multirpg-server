package main

import "net"

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
