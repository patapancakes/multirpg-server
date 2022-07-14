package main

type Message struct {
	sender *Client
	data []byte
}

func (m *Message) process() {

}
