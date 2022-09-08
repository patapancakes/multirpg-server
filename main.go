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
	"flag"
	"fmt"
	"os"
)

func main() {
	server := &Server{
		charSets:  getCharSetList(),

		rooms:     make(map[uint16]*Room),
		clientIds: make(map[uint16]bool),
	}

	if err := server.start(readFlags()); err != nil {
		fmt.Printf("server error: %s\n", err)
		os.Exit(1)
	}
}

func readFlags() (*string, *int) {
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.Int("port", 22888, "Port to listen on")
	flag.Parse()

	return host, port
}
