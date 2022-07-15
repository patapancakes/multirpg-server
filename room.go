package main

type Room struct {
	server *Server
	id uint16

	clients map[*Client]bool
}

func (s *Server) createRoom(id uint16) *Room {
	return &Room{
		server: s,
		id: id,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) broadcast(data []byte, sender *Client) {
	// Packet broadcasting is disabled while on the title screen
	if r.id == 0 {
		return
	}

	for client := range r.clients {
		if client == sender {
			continue
		}

		client.conn.Write(data)
	}
}
