package main

/*
multirpg-server
https://github.com/Gamizard/multirpg-server

Copyright (C) 2022 azarashi <azarashi@majestaria.fun>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"net"
)

type Server struct {
	charSets  []string

	rooms     map[uint16]*Room
	clientIds map[uint16]bool
}

func (s *Server) start(host *string, port *int) error {
	fmt.Println("Starting server on " + *host + ":" + fmt.Sprint(*port))

	// Room 0 is the room clients are put in when they first connect
	// Clients are expected to send a switch room packet to join a game room
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
		id:   s.getFreeId(),
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
