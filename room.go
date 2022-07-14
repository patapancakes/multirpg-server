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

func (r *Room) getFreeId() uint16 {
	for i := uint16(1); i < 0xFFFF; i++ {
		if _, ok := r.clients[i]; !ok {
			return i
		}
	}

	return 0
}
