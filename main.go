package main

func main() {
	server := &Server{
		rooms: make(map[uint16]*Room),
	}

	for _, mapID := range getMapList() {
		server.rooms[mapID] = &Room{
			id: mapID,
			clients: make(map[uint16]*Client),
		}
	}

	server.start()
}
