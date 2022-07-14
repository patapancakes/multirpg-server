package main

type Room struct {
	server *Server
	id uint16
	clients map[uint16]*Client
}

func (s *Server) createRoom(id uint16) *Room {
	return &Room{
		server: s,
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

func (r *Room) broadcast(data []byte) {
	for _, client := range r.clients {
		client.conn.Write(data)
	}
}
