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
	rooms     map[uint16]*Room
	clientIds map[uint16]bool
}

func (s *Server) start(host *string, port *int) error {
	fmt.Printf("Starting server on %s:%d\n", *host, *port)

	// Room 0 is the room clients are put in when they first connect
	// Clients are expected to send a switch room packet to join a game room
	s.rooms[0] = s.createRoom(0)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		return err
	}

	// Listen for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	client := &Client{
		conn: conn,
		id:   s.getFreeId(),
	}

	s.clientIds[client.id] = true
	
	fmt.Printf("Connection from %s (client %d)\n", conn.RemoteAddr().String(), client.id)
	
	client.joinRoom(0)

	client.listen()

	client.disconnect()
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
