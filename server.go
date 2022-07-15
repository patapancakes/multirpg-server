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

	if maps, err := getMapList(); err != nil {
		return err
	} else {
		for _, mapId := range maps {
			s.rooms[mapId] = s.createRoom(mapId)
		}
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
	client.leaveRoom()
	delete(s.clientIds, client.id)
}

func (s *Server) getFreeId() uint16 {
	for i := uint16(0); i < 0xFFFF; i++ {
		if _, ok := s.clientIds[i]; !ok {
			return i
		}
	}

	// This should never happen, if it does then somehow all ids are being used
	return 0
}
