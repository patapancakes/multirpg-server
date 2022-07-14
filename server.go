package main

import (
	"fmt"
	"net"
)

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

}
