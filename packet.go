package main

type Packet struct {
	sender *Client
	data []byte
}

func (m *Packet) process() {

}
