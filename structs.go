package main

import "net"

type Server struct {
	rooms map[uint16]*Room
}

type Room struct {
	id uint16
	clients map[uint16]*Client
}

type Client struct {
	conn net.Conn
	room *Room
	id uint16
}

type Message struct {
	sender *Client
	data []byte
}
