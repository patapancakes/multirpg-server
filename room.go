package main

type Room struct {
	id uint16
	clients map[uint16]*Client
}

func createRoom(id uint16) *Room {
	return &Room{
		id: id,
		clients: make(map[uint16]*Client),
	}
}
