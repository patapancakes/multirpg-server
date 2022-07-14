package main

import "net"

type Client struct {
	conn net.Conn
	room *Room
	id uint16
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

		go packet.process()
	}
}
