package main

type Server struct {
	rooms map[uint16]*Room
}

type Room struct {
	id uint16
	clients map[uint16]*Client
}

type Client struct {
	id uint16
	room *Room
}
