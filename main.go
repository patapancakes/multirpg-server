package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	server := &Server{
		rooms: make(map[uint16]*Room),
	}

	server.rooms[0] = &Room{
		id: 0,
		clients: make(map[uint16]*Client),
	}

	for _, mapID := range getMapList() {
		server.rooms[mapID] = &Room{
			id: mapID,
			clients: make(map[uint16]*Client),
		}
	}

	if err := server.start(readFlags()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readFlags() (*string, *int) {
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.Int("port", 22888, "Port to listen on")
	flag.Parse()

	return host, port
}
