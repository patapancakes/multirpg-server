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
	charSets map[string]bool

	rooms     map[uint16]*Room
	clientIds map[uint16]bool
}

func (s *Server) start(host *string, port *int) error {
	fmt.Printf("Starting server on %s:%d\n", *host, *port)

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

	// Listen for incoming connections
	for {
		conn, err := server.Accept()
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
		room: s.rooms[0],
	}

	s.clientIds[client.id] = true
	s.rooms[0].clients[client] = true

	fmt.Printf("Connection from %s (client %d)\n", conn.RemoteAddr().String(), client.id)

	client.listen()

	client.leaveRoom()

	// Release client id
	delete(s.clientIds, client.id)

	if err := conn.Close(); err != nil {
		fmt.Printf("Connection from %s (client %d) failed to close: %s\n", conn.RemoteAddr().String(), client.id, err)
		return
	}

	fmt.Printf("Connection from %s (client %d) closed\n", conn.RemoteAddr().String(), client.id)
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
