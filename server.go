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

package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	lobbies sync.Map
}

func (s *Server) start(host *string, port *int) error {
	fmt.Printf("Starting server on %s:%d\n", *host, *port)

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
		conn:    conn,
		send:    make(chan []byte),
		receive: make(chan []byte),
		server:  s,
	}

	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())

	go client.packetWriter()
	go client.packetReader()

	client.listen()

	client.disconnect()
}
