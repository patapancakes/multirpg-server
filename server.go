package main

import (
	"fmt"
	"net"
)

type Server struct {
	rooms map[uint16]*Room
}

func (s *Server) start(host *string, port *int) error {
	fmt.Println("Starting server on " + *host + ":" + fmt.Sprint(*port))

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		return err
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			return err
		}

		fmt.Println("Connection from " + conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := &Client{
		conn: conn,
		room: s.rooms[0],
		id: s.rooms[0].getFreeId(),
	}

	s.rooms[0].clients[client.id] = client

	client.listen()

	fmt.Println("Connection from " + conn.RemoteAddr().String() + " closed")
	delete(s.rooms[client.room.id].clients, client.id)
}
