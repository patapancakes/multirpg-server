package main

import (
	"fmt"
	"net"
)

type Server struct {
	rooms map[uint16]*Room
	clientIds map[uint16]bool
}

func (s *Server) start(host *string, port *int) error {
	fmt.Println("Starting server on " + *host + ":" + fmt.Sprint(*port))

	s.rooms[0] = s.createRoom(0)

	for _, mapID := range getMapList() {
		s.rooms[mapID] = s.createRoom(mapID)
	}

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
		id: s.getFreeId(),
	}

	s.rooms[0].clients[client] = true

	client.listen()

	fmt.Println("Connection from " + conn.RemoteAddr().String() + " closed")
	client.handleDisconnect()
	delete(s.clientIds, client.id)
}

func (s *Server) getFreeId() uint16 {
	for i := uint16(1); i < 0xFFFF; i++ {
		if _, ok := s.clientIds[i]; !ok {
			return i
		}
	}

	return 0
}
