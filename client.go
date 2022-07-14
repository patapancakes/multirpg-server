package main

func (c *Client) listen() {
	for	{
		buf := make([]byte, 1024)

		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}

		message := &Message{
			sender: c,
			data: buf[:n],
		}

		go message.process()
	}
}
