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
	"os"
	"path/filepath"
	"strconv"
)

func getMapList() ([]uint16, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var maps []uint16
	for _, file := range files {
		if len(file.Name()) == 11 && file.Name()[7:] == ".lmu" {
			id, err := strconv.Atoi(file.Name()[3:7])
			if err != nil {
				fmt.Printf("Unable to get map id from filename (%s), skipping...\n", file.Name())
				continue
			}

			maps = append(maps, uint16(id))
		}
	}

	if len(maps) < 1 {
		fmt.Print("No maps were found\nMultiplayer map changes will not work. Make sure you're running multirpg-server from the game data folder and have permission to read files\n")
	}

	return maps, nil
}

func getCharSetList() map[string]bool {
	files, err := os.ReadDir("CharSet")
	if err != nil {
		fmt.Printf("%s\nMultiplayer sprite changes will not work. Make sure you're running multirpg-server from the game data folder and have permission to read files\n", err)
	}

	charSets := make(map[string]bool)
	for _, file := range files {
		charSets[file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]] = true
	}

	return charSets
}

/*func getSoundList() map[string]bool {
	files, err := os.ReadDir("Sound")
	if err != nil {
		fmt.Printf("%s\nMultiplayer sounds will not work. Make sure you're running multirpg-server from the game data folder and have permission to read files\n", err)
	}

	sounds := make(map[string]bool)
	for _, file := range files {
		sounds[file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]] = true
	}

	return sounds
}

func getSystemList() map[string]bool {
	files, err := os.ReadDir("System")
	if err != nil {
		fmt.Printf("%s\nMultiplayer system changes will not work. Make sure you're running multirpg-server from the game data folder and have permission to read files\n", err)
	}

	systems := make(map[string]bool)
	for _, file := range files {
		systems[file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]] = true
	}

	return systems
}*/

func (s *Server) isValidCharSet(charSet string) bool {
	return s.charSets[charSet]
}
